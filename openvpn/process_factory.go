/*
 * go-openvpn -- Go gettable library for wrapping Openvpn functionality in go way.
 *
 * Copyright (C) 2020 The "MysteriumNetwork/go-openvpn" Authors..
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License Version 3
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program in the COPYING file.
 * If not, see <http://www.gnu.org/licenses/>.
 */

package openvpn

import (
	"os/exec"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/config"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/tunnel"
)

// CreateNewProcess creates new openvpn process with given config params
func CreateNewProcess(openvpnBinary string, config *config.GenericConfig, middlewares ...management.Middleware) *OpenvpnProcess {
	tunnelSetup := tunnel.NewTunnelSetup()
	execCommand := func(arg ...string) *exec.Cmd {
		return exec.Command(openvpnBinary, arg...)
	}
	return newProcess(tunnelSetup, config, execCommand, middlewares...)
}
