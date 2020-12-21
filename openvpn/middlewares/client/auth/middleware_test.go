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

package auth

import (
	"testing"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
	"github.com/stretchr/testify/assert"
)

func auth() (string, string, error) {
	return "testuser", "testpassword", nil
}

func Test_Factory(t *testing.T) {
	middleware := NewMiddleware(auth)
	assert.NotNil(t, middleware)
}

func Test_ConsumeLineSkips(t *testing.T) {
	var tests = []struct {
		line string
	}{
		{">SOME_LINE_DELIVERED"},
		{">ANOTHER_LINE_DELIVERED"},
		{">PASSWORD"},
	}
	middleware := NewMiddleware(auth)

	for _, test := range tests {
		consumed, err := middleware.ConsumeLine(test.line)
		assert.NoError(t, err, test.line)
		assert.False(t, consumed, test.line)
	}
}

func Test_ConsumeLineTakes(t *testing.T) {
	passwordRequest := ">PASSWORD:Need 'Auth' username/password"

	middleware := NewMiddleware(auth)
	mockCmdWriter := &management.MockConnection{}
	middleware.Start(mockCmdWriter)

	consumed, err := middleware.ConsumeLine(passwordRequest)
	assert.NoError(t, err)
	assert.True(t, consumed)
	assert.Equal(t,
		mockCmdWriter.WrittenLines,
		[]string{
			"password 'Auth' testpassword",
			"username 'Auth' testuser",
		},
	)
}
