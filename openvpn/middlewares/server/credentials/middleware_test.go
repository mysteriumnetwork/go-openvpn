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

package credentials

import (
	"testing"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/middlewares/server"
	"github.com/stretchr/testify/assert"
)

type fakeValidator struct {
	called bool
}

func (f *fakeValidator) authenticateClient(clientID int, username, password string) (bool, error) {
	f.called = true
	return username == "username1" && password == "12341234", nil
}

func Test_Factory(t *testing.T) {
	fas := fakeValidator{}
	middleware := NewMiddleware(fas.authenticateClient)
	assert.NotNil(t, middleware)
}

func Test_ConsumeLineAuthTrueChecker(t *testing.T) {
	fas := fakeValidator{}

	mockConnection := &management.MockConnection{}
	middleware := NewMiddleware(fas.authenticateClient)
	middleware.Start(mockConnection)

	middleware.handleClientEvent(server.ClientEvent{
		EventType: server.Connect,
		ClientID:  3,
		ClientKey: 4,
		Env: map[string]string{
			"username": "username1",
			"password": "12341234",
		},
	})
	assert.Equal(t, "client-auth-nt 3 4", mockConnection.LastLine)
}

func Test_ConsumeLineAuthFalseChecker(t *testing.T) {
	fas := fakeValidator{}

	mockConnection := &management.MockConnection{}
	middleware := NewMiddleware(fas.authenticateClient)
	middleware.Start(mockConnection)

	middleware.handleClientEvent(server.ClientEvent{
		EventType: server.Connect,
		ClientID:  3,
		ClientKey: 4,
	})
	assert.Equal(t, "client-deny 3 4 missing username or password", mockConnection.LastLine)

	middleware.handleClientEvent(server.ClientEvent{
		EventType: server.Connect,
		ClientID:  3,
		ClientKey: 4,
		Env: map[string]string{
			"username": "username1",
			"password": "wrong",
		},
	})
	assert.Equal(t, "client-deny 3 4 wrong username or password", mockConnection.LastLine)
}
