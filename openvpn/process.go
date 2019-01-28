/*
 * Copyright (C) 2018 The "MysteriumNetwork/go-openvpn" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package openvpn

import (
	"errors"
	"os/exec"
	"sync"
	"time"

	log "github.com/cihub/seelog"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/config"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/tunnel"
)

const openvpnManagementLogPrefix = "[client-management] "
const openvpnProcessLogPrefix = "[openvpn-process] "

// OpenvpnProcess represents an openvpn process manager
type OpenvpnProcess struct {
	config              *config.GenericConfig
	tunnelSetup         tunnel.Setup
	management          *management.Management
	cmd                 *CmdWrapper
	lastSessionShutdown chan bool
}

func newProcess(
	tunnelSetup tunnel.Setup,
	config *config.GenericConfig,
	execCommand func(arg ...string) *exec.Cmd,
	lastSessionShutdown chan bool,
	middlewares ...management.Middleware,
) *OpenvpnProcess {
	return &OpenvpnProcess{
		tunnelSetup:         tunnelSetup,
		config:              config,
		management:          management.NewManagement(management.LocalhostOnRandomPort, openvpnManagementLogPrefix, middlewares...),
		cmd:                 NewCmdWrapper(openvpnProcessLogPrefix, execCommand),
		lastSessionShutdown: lastSessionShutdown,
	}
}

// Start starts the openvpn process
func (openvpn *OpenvpnProcess) Start() error {
	if err := openvpn.tunnelSetup.Setup(openvpn.config); err != nil {
		return err
	}

	err := openvpn.management.WaitForConnection()
	if err != nil {
		openvpn.tunnelSetup.Stop()
		return err
	}

	addr := openvpn.management.BoundAddress
	openvpn.config.SetManagementAddress(addr.IP, addr.Port)

	// Fetch the current arguments
	arguments, err := (*openvpn.config).ToArguments()
	if err != nil {
		log.Info(openvpnManagementLogPrefix, "stopping management on argument error")
		openvpn.management.Stop()
		openvpn.tunnelSetup.Stop()
		return err
	}

	//nil returned from process.Start doesn't guarantee that openvpn itself initialized correctly and accepted all arguments
	//it simply means that OS started process with specified args
	err = openvpn.cmd.Start(arguments)
	if err != nil {
		log.Info(openvpnManagementLogPrefix, "stopping management on openvpn start error")
		openvpn.management.Stop()
		openvpn.tunnelSetup.Stop()
		return err
	}

	select {
	case connAccepted := <-openvpn.management.Connected:
		if connAccepted {
			return nil
		}
		return errors.New("management failed to accept connection")
	case exitError := <-openvpn.cmd.CmdExitError:
		log.Info(openvpnManagementLogPrefix, "stopping management on previous openvpn exit error: ", exitError)
		openvpn.management.Stop()
		openvpn.tunnelSetup.Stop()
		if exitError != nil {
			return exitError
		}
		return errors.New("openvpn process died too early")
	case <-time.After(2 * time.Second):
		return errors.New("management connection wait timeout")
	}
}

// Wait waits for the openvpn process to complete
func (openvpn *OpenvpnProcess) Wait() error {
	for {
		select {
		case lastSessionShutdown := <-openvpn.lastSessionShutdown:
			if lastSessionShutdown {
				log.Info(openvpnManagementLogPrefix, "exiting openvpn process since last session has been closed")
				openvpn.Stop()
				return nil
			}
		case exitError := <-openvpn.cmd.CmdExitError:
			return exitError
		}
	}
	return nil
}

// Stop stops the openvpn process
func (openvpn *OpenvpnProcess) Stop() {
	waiter := sync.WaitGroup{}
	//TODO which to signal for close first ?
	//if we stop process before management, managemnt won't have a chance to send any commands from middlewares on stop
	//if we stop management first - it will miss important EXITING state from process
	waiter.Add(1)
	go func() {
		defer waiter.Done()
		openvpn.cmd.Stop()
	}()

	waiter.Add(1)
	go func() {
		defer waiter.Done()
		log.Info("stopping management on openvpn process stop after defer waiter")
		openvpn.management.Stop()
	}()
	waiter.Wait()

	openvpn.tunnelSetup.Stop()
}
