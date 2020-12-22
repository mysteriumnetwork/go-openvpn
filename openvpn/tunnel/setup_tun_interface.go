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

import (
	"errors"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/config"
)

const tunLogPrefix = "[linux tun service] "

// ErrNoFreeTunDevice is thrown when no free tun device is available on system
var ErrNoFreeTunDevice = errors.New("no free tun device found")

// Setup represents the operations required for a tunnel setup
type Setup interface {
	Setup(config *config.GenericConfig) error
	Stop()
	DeviceName() string
}
