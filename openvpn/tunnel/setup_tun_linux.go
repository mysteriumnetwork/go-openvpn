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

package tunnel

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/config"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/log"
)

// NewTunnelSetup returns a new tunnel setup for linux
func NewTunnelSetup() Setup {
	return &LinuxTunDeviceManager{}
}

// LinuxTunDeviceManager represents the tun manager for linux
type LinuxTunDeviceManager struct {
	scriptSetup string

	// runtime variables
	device tunDevice
}

// tunDevice represents tun device structure
type tunDevice struct {
	Name string
}

// Setup sets the tunel up
func (service *LinuxTunDeviceManager) Setup(configuration *config.GenericConfig) error {
	configuration.SetScriptParam("iproute", config.SimplePath("nonpriv-ip"))
	service.scriptSetup = configuration.GetFullScriptPath(config.SimplePath("prepare-env.sh"))

	err := service.createDeviceNode()
	if err != nil {
		return errors.Wrap(err, "failed to create device node")
	}

	device, err := service.getNextFreeTunDevice()
	if err != nil {
		return err
	}

	service.device = device
	configuration.SetPersistTun()
	configuration.SetDevice(device.Name)
	return nil
}

// Stop destroys tunnel device
func (service *LinuxTunDeviceManager) Stop() {
	var err error
	var exists bool

	if exists, err = service.deviceExists(service.device.Name); err != nil {
		log.Info(tunLogPrefix, err)
	}

	if exists {
		service.deleteDevice(service.device)
	}
}

// DeviceName returns tunnel device name
func (service *LinuxTunDeviceManager) DeviceName() string {
	return service.device.Name
}

func (service *LinuxTunDeviceManager) createTunDevice(deviceName string) (err error) {
	cmd := exec.Command("sudo", "ip", "tuntap", "add", "dev", deviceName, "mode", "tun")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Warn("Failed to add tun device:", cmd.Args, "Returned exit error:", err, "Cmd output:", string(output))
		// we should not proceed without tun device
		return err
	}

	log.Info(tunLogPrefix, deviceName+"device created")
	return nil
}

func deviceUsed(deviceName string) (used bool, err error) {
	contents, err := ioutil.ReadFile("/sys/class/net/" + deviceName + "/carrier")
	if err != nil {
		return false, err
	}

	value, err := strconv.Atoi(strings.TrimSuffix(string(contents), "\n"))
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

func (service *LinuxTunDeviceManager) deleteDevice(device tunDevice) {
	// Cleaning here as much as possible, if device deletion failed we at least unassigned IP-addresses.
	cmd := exec.Command("sudo", "ip", "addr", "flush", "dev", device.Name)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Warn("Failed to flush tun device:", cmd.Args, "Returned exit error:", err, "Cmd output:", string(output))
	}

	cmd = exec.Command("sudo", "ip", "tuntap", "delete", "dev", device.Name, "mode", "tun")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Warn("Failed to remove tun device:", cmd.Args, "Returned exit error:", err, "Cmd output:", string(output))
	} else {
		log.Info(tunLogPrefix, device.Name, "device removed")
	}
}

// getNextFreeTunDevice returns first free tun device on system
func (service *LinuxTunDeviceManager) getNextFreeTunDevice() (tun tunDevice, err error) {
	// search only among first 10 tun devices
	for i := 0; i <= 10; i++ {
		tunName := "tun" + strconv.Itoa(i)
		tunFile := "/sys/class/net/" + tunName
		if _, err := os.Stat(tunFile); err == nil {
			used, err := deviceUsed(tunName)
			if err != nil {
				return tunDevice{}, errors.Wrap(err, "failed to check if device is used")
			}
			if !used {
				log.Debug("Tunnel exists, but not used, reusing:" + tunFile)
				return tunDevice{tunName}, nil
			}
			log.Debug("Tunnel exists and is taken:" + tunFile)
		} else if os.IsNotExist(err) {
			log.Debug("Tunnel does not exists, creating:" + tunFile)
			err := service.createTunDevice(tunName)
			if err != nil {
				return tunDevice{}, errors.Wrap(err, "failed to create a tunnel: "+tunFile)
			}
			return tunDevice{tunName}, nil
		} else if err != nil {
			log.Error("Failed to check if tunnel device exists:", err)
		}
	}

	return tun, ErrNoFreeTunDevice
}

func (service *LinuxTunDeviceManager) deviceExists(tunName string) (exists bool, err error) {
	tunFile := "/sys/class/net/" + tunName
	if _, err := os.Stat(tunFile); err == nil {
		return true, nil
	}
	return false, err
}

func (service *LinuxTunDeviceManager) createDeviceNode() error {
	if _, err := os.Stat("/dev/net/tun"); err == nil {
		// device node already exists
		return nil
	}

	cmd := exec.Command("sudo", service.scriptSetup)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Warn("Failed to execute tun script:", cmd.Args, "Returned exit error:", err, "Cmd output:", string(output))
		return err
	}

	log.Info(tunLogPrefix, "/dev/net/tun device node created")
	return nil
}
