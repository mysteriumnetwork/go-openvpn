/*
 * go-openvpn -- Go gettable library for wrapping Openvpn functionality in go way.
 *
 * Copyright (C) 2020 BlockDev AG.
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
