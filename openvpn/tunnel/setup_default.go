/*
 * go-openvpn -- Go gettable library for wrapping Openvpn functionality in go way.
 *
 * Copyright (C) 2020 BlockDev AG.
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

package tunnel

import "github.com/mysteriumnetwork/go-openvpn/openvpn/config"

// DefaultSetup represents a default tunnel setup - aka it sets the tun in configuration
type DefaultSetup struct {
}

// Setup implements the setup method for tunnel interface
func (ds *DefaultSetup) Setup(config *config.GenericConfig) error {
	config.SetDevice("tun")
	return nil
}

// Stop implements the stop method for tunnel interface
func (ds *DefaultSetup) Stop() {

}

// DeviceName returns tunnel device name
func (ds *DefaultSetup) DeviceName() string {
	return "tun"
}
