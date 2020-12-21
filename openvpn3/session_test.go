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

package openvpn3

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestSessionStartStopDoesNotPanic(t *testing.T) {
	session := NewSession(Config{}, UserCredentials{}, &fmtLogger{})
	session.Start()
	session.Stop()
}

func TestSessionInitFailsForInvalidProfile(t *testing.T) {
	session := NewSession(Config{}, UserCredentials{}, &fmtLogger{})
	session.Start()
	err := session.Wait()
	assert.Equal(t, ErrInitFailed, err)
}

func TestSessionConnectFailsForInvalidRemote(t *testing.T) {
	session := NewSession(NewConfig("remote localhost 1111"), UserCredentials{}, &fmtLogger{})
	session.Start()
	err := session.Wait()
	assert.Equal(t, ErrConnectFailed, err)
}

type fmtLogger struct {
}

func (l *fmtLogger) Log(text string) {
	fmt.Println(text)
}
