/*
 * Copyright (C) 2019 The "MysteriumNetwork/go-openvpn" Authors.
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

package bytecount

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
)

var rule = regexp.MustCompile("^>BYTECOUNT_CLI:([0-9]*),([0-9]*),([0-9]*)$")

// SessionByteChangeHandler is the callback we call with the session byte count
type SessionByteChangeHandler func(SessionByteCount)

// SessionByteCount represents
type SessionByteCount struct {
	ClientID          int
	BytesIn, BytesOut uint64
}

// Middleware reports the different session byte counts
type Middleware struct {
	handler        SessionByteChangeHandler
	updateInterval int
}

// NewMiddleware returns a new instance of the middleware
func NewMiddleware(h SessionByteChangeHandler, updateInterval int) *Middleware {
	return &Middleware{
		handler:        h,
		updateInterval: updateInterval,
	}
}

// Start starts the middleware
func (m *Middleware) Start(cw management.CommandWriter) error {
	_, err := cw.SingleLineCommand("bytecount %v", m.updateInterval)
	return err
}

// Stop stops the middleware
func (m *Middleware) Stop(cw management.CommandWriter) error {
	_, err := cw.SingleLineCommand("bytecount %v", 0)
	return err
}

// ConsumeLine handles the given openvpn management line
func (m *Middleware) ConsumeLine(line string) (consumed bool, err error) {
	if !rule.MatchString(line) {
		return false, nil
	}

	match := rule.FindStringSubmatch(line)

	if len(match) != 4 {
		return false, fmt.Errorf("wrong match length for %q. got len = %v, expected %v", line, len(match), 4)
	}

	clientID, err := strconv.Atoi(match[1])
	if err != nil {
		return false, fmt.Errorf("could not parse clientID from match[1]: %q", match[1])
	}

	bytesIn, err := strconv.ParseUint(match[2], 10, 64)
	if err != nil {
		return false, fmt.Errorf("could not parse clientID from match[2]: %q", match[2])
	}

	bytesOut, err := strconv.ParseUint(match[3], 10, 64)
	if err != nil {
		return false, fmt.Errorf("could not parse clientID from match[3]: %q", match[3])
	}

	m.handler(SessionByteCount{
		ClientID: clientID,
		BytesIn:  bytesIn,
		BytesOut: bytesOut,
	})

	return true, nil
}
