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
