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

import (
	"fmt"
)

// MockConnection is mock openvpn management interface used for middleware testing
type MockConnection struct {
	WrittenLines      []string
	LastLine          string
	CommandResult     string
	MultilineResponse []string
}

// SingleLineCommand sends command to mocked connection and expects single line as command output (error or success)
func (conn *MockConnection) SingleLineCommand(format string, args ...interface{}) (string, error) {
	conn.LastLine = fmt.Sprintf(format, args...)
	conn.WrittenLines = append(conn.WrittenLines, conn.LastLine)
	return conn.CommandResult, nil
}

// MultiLineCommand sends command to mocked connection and expects multiple line command response with END marker
func (conn *MockConnection) MultiLineCommand(format string, args ...interface{}) (string, []string, error) {
	_, _ = conn.SingleLineCommand(format, args...)
	return conn.CommandResult, conn.MultilineResponse, nil
}
