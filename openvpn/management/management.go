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

package management

// Management packages contains all functionality related to openvpn management interface
// See https://openvpn.net/index.php/open-source/documentation/miscellaneous/79-management-interface.html

// CommandWriter represents openvpn management interface abstraction for middlewares to be able to send commands to openvpn process
type CommandWriter interface {
	SingleLineCommand(template string, args ...interface{}) (string, error)
	MultiLineCommand(template string, args ...interface{}) (string, []string, error)
}

// Middleware used to control openvpn process through management interface
// It's guaranteed that ConsumeLine callback will be called AFTER Start callback is finished
// CommandWriter passed on Stop callback can be already closed - expect errors when sending commands
// For efficiency and simplicity purposes ConsumeLine for each middleware is called from the same goroutine which
// consumes events from channel - avoid long running operations at all costs
type Middleware interface {
	Start(CommandWriter) error
	Stop(CommandWriter) error
	ConsumeLine(line string) (consumed bool, err error)
}
