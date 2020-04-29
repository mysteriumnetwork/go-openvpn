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

package auth

import (
	"testing"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/middlewares/server"
	"github.com/stretchr/testify/assert"
)

func Test_ConsumeLineSkips(t *testing.T) {
	var tests = []struct {
		line string
	}{
		{">SOME_LINE_TO_BE_DELIVERED"},
		{">ANOTHER_LINE_TO_BE_DELIVERED"},
		{">PASSWORD"},
		{">USERNAME"},
	}
	middleware := NewMiddleware()

	for _, test := range tests {
		consumed, err := middleware.ConsumeLine(test.line)
		assert.NoError(t, err, test.line)
		assert.False(t, consumed, test.line)
	}
}

func Test_ConsumeLineTakes(t *testing.T) {
	var tests = []struct {
		line string
	}{
		{">CLIENT:REAUTH,0,0"},
		{">CLIENT:CONNECT,0,0"},
		{">CLIENT:ENV,password=12341234"},
		{">CLIENT:ENV,username=username"},
	}

	middleware := NewMiddleware()
	mockConnection := &management.MockConnection{}
	middleware.Start(mockConnection)

	for _, test := range tests {
		consumed, err := middleware.ConsumeLine(test.line)
		assert.NoError(t, err, test.line)
		assert.True(t, consumed, test.line)
	}
}

func Test_ConsumeLineShouldNotTriggerClientState(t *testing.T) {
	var tests = []struct {
		line string
	}{
		{">CLIENT:ENV,password=12341234"},
		{">CLIENT:ENV,username=username"},
		{">CLIENT:REAUTH,0,0"},
		{">CLIENT:CONNECT,0,0"},
	}

	for _, test := range tests {
		var receivedEvent *server.ClientEvent
		middleware := NewMiddleware(func(e server.ClientEvent) {
			receivedEvent = &e
		})

		mockConnection := &management.MockConnection{}
		middleware.Start(mockConnection)

		consumed, err := middleware.ConsumeLine(test.line)
		assert.NoError(t, err, test.line)
		assert.True(t, consumed, test.line)
		assert.Nil(t, receivedEvent)
	}
}

func Test_ConsumeLineShouldTriggerClientStateAfterReceivingEnvironment(t *testing.T) {
	var tests = []struct {
		line string
	}{
		{">CLIENT:CONNECT,1,2"},
		{">CLIENT:ENV,password=12341234"},
		{">CLIENT:ENV,username=username1"},
		{">CLIENT:ENV,END"},
	}

	var receivedEvent server.ClientEvent
	middleware := NewMiddleware(func(e server.ClientEvent) {
		receivedEvent = e
	})

	mockConnection := &management.MockConnection{}
	middleware.Start(mockConnection)

	for _, test := range tests {
		consumed, err := middleware.ConsumeLine(test.line)
		assert.NoError(t, err, test.line)
		assert.True(t, consumed, test.line)
	}
	assert.Equal(
		t,
		server.ClientEvent{
			EventType: server.Connect,
			ClientID:  1,
			ClientKey: 2,
			Env: map[string]string{
				"username": "username1",
				"password": "12341234",
			},
		},
		receivedEvent,
	)
}

func Test_ConsumeLinesAcceptsClientIdsAntKeysWithSeveralDigits(t *testing.T) {
	var tests = []string{
		">CLIENT:CONNECT,115,23",
		">CLIENT:ENV,END",
	}

	var receivedEvent server.ClientEvent
	middleware := NewMiddleware(func(e server.ClientEvent) {
		receivedEvent = e
	})

	mockConnection := &management.MockConnection{}
	middleware.Start(mockConnection)

	for _, testLine := range tests {
		consumed, err := middleware.ConsumeLine(testLine)
		assert.NoError(t, err, testLine)
		assert.Equal(t, true, consumed, testLine)
	}

	assert.Equal(
		t,
		server.ClientEvent{
			EventType: server.Connect,
			ClientID:  115,
			ClientKey: 23,
			Env:       map[string]string{},
		},
		receivedEvent,
	)
}
