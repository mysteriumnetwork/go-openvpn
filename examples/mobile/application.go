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

package mobile

import (
	"fmt"
	"strings"

	"github.com/mysteriumnetwork/go-openvpn/openvpn3"
)

type callbacks interface {
	openvpn3.Logger
	openvpn3.EventConsumer
	openvpn3.StatsConsumer
}

type loggingCallbacks struct {
}

func (lc *loggingCallbacks) Log(text string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		fmt.Println("Openvpn log >>", line)
	}
}

func (lc *loggingCallbacks) OnEvent(event openvpn3.Event) {
	fmt.Printf("Openvpn event >> %+v\n", event)
}

func (lc *loggingCallbacks) OnStats(stats openvpn3.Statistics) {
	fmt.Printf("Openvpn stats >> %+v\n", stats)
}

var _ callbacks = &loggingCallbacks{}

// StdoutLogger represents the stdout logger callback
type StdoutLogger func(text string)

// Log logs the given string to stdout logger
func (lc StdoutLogger) Log(text string) {
	lc(text)
}

// StartSession starts the openvpn3 session for ios
func StartSession() {

	var logger StdoutLogger = func(text string) {
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			fmt.Println("Library check >>", line)
		}
	}
	openvpn3.SelfCheck(logger)

	session := openvpn3.NewMobileSession(
		openvpn3.NewConfig(""),
		openvpn3.UserCredentials{},
		&loggingCallbacks{},
		&openvpn3.NoOpTunnelSetup{},
	)

	session.Start()
	err := session.Wait()
	if err != nil {
		fmt.Println("Openvpn3 error: ", err)
	} else {
		fmt.Println("Graceful exit")
	}

}
