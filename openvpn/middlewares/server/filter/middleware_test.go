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

package filter

import (
	"testing"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/middlewares/server"
	"github.com/stretchr/testify/assert"
)

const (
	emptyFilter = `client-pf 0
[CLIENTS DROP]
[SUBNETS ACCEPT]
[END]
END
`

	allowFilter = `client-pf 0
[CLIENTS DROP]
[SUBNETS ACCEPT]
+1.1.1.1/32
+2.2.2.0/24
[END]
END
`

	blockFilter = `client-pf 0
[CLIENTS DROP]
[SUBNETS ACCEPT]
-1.1.1.1/32
-2.2.2.0/24
[END]
END
`

	bothFilter = `client-pf 0
[CLIENTS DROP]
[SUBNETS ACCEPT]
+1.1.1.1/32
+2.2.2.0/24
-3.3.3.3/32
-4.4.4.0/24
[END]
END
`
)

func Test_EmptyPacketFilterForEmptySubnets(t *testing.T) {
	mockConnection := &management.MockConnection{}
	middleware := NewMiddleware(nil, nil)
	middleware.Start(mockConnection)

	middleware.handleClientEvent(server.ClientEvent{EventType: server.Connect, ClientID: 0})
	assert.Equal(t, emptyFilter, mockConnection.WrittenLines[0])
}

func Test_AllowPacketFilterForAllowSubnets(t *testing.T) {
	middleware := NewMiddleware([]string{"1.1.1.1/32", "2.2.2.0/24"}, nil)
	mockConnection := &management.MockConnection{}
	middleware.Start(mockConnection)

	middleware.handleClientEvent(server.ClientEvent{EventType: server.Connect, ClientID: 0})
	assert.Equal(t, allowFilter, mockConnection.WrittenLines[0])
}

func Test_BlockPacketFilterForBlockSubnets(t *testing.T) {
	middleware := NewMiddleware(nil, []string{"1.1.1.1/32", "2.2.2.0/24"})
	mockConnection := &management.MockConnection{}
	middleware.Start(mockConnection)

	middleware.handleClientEvent(server.ClientEvent{EventType: server.Connect, ClientID: 0})
	assert.Equal(t, blockFilter, mockConnection.WrittenLines[0])
}

func Test_BothPacketFilterForBothSubnets(t *testing.T) {
	middleware := NewMiddleware([]string{"1.1.1.1/32", "2.2.2.0/24"}, []string{"3.3.3.3/32", "4.4.4.0/24"})
	mockConnection := &management.MockConnection{}
	middleware.Start(mockConnection)

	middleware.handleClientEvent(server.ClientEvent{EventType: server.Connect, ClientID: 0})
	assert.Equal(t, bothFilter, mockConnection.WrittenLines[0])
}
