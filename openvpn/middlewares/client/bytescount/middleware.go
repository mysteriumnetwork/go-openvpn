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

package bytescount

import (
	"regexp"
	"strconv"
	"time"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
)

// SessionStatsHandler is invoked when middleware receives statistics
type SessionStatsHandler func(Bytecount) error

// Bytecount represents the bytecount response
type Bytecount struct {
	BytesIn  uint64
	BytesOut uint64
}

const byteCountCommandTemplate = "bytecount %d"

var rule = regexp.MustCompile("^>BYTECOUNT:(.*),(.*)$")

type middleware struct {
	sessionStatsHandler SessionStatsHandler
	interval            time.Duration
}

// NewMiddleware returns new bytescount middleware
func NewMiddleware(sessionStatsHandler SessionStatsHandler, interval time.Duration) management.Middleware {
	return &middleware{
		sessionStatsHandler: sessionStatsHandler,
		interval:            interval,
	}
}

func (middleware *middleware) Start(commandWriter management.CommandWriter) error {
	_, err := commandWriter.SingleLineCommand(byteCountCommandTemplate, int(middleware.interval.Seconds()))
	return err
}

func (middleware *middleware) Stop(commandWriter management.CommandWriter) error {
	_, err := commandWriter.SingleLineCommand(byteCountCommandTemplate, 0)
	return err
}

func (middleware *middleware) ConsumeLine(line string) (consumed bool, err error) {
	match := rule.FindStringSubmatch(line)
	if consumed = len(match) > 2; !consumed {
		return
	}

	bytesIn, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return
	}

	bytesOut, err := strconv.ParseUint(match[2], 10, 64)
	if err != nil {
		return
	}

	err = middleware.sessionStatsHandler(Bytecount{BytesIn: bytesIn, BytesOut: bytesOut})

	return
}
