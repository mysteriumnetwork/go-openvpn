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

package server

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientEventIsParsed(t *testing.T) {
	var testData = []struct {
		testLine  string
		event     ClientEventType
		eventData string
		err       error
	}{
		{"CONNECT,1,1", Connect, "1,1", nil},
		{"REAUTH,2,2", Reauth, "2,2", nil},
		{"ENV,abc=123", Env, "abc=123", nil},
		{"ESTABLISHED,1", Established, "1", nil},
		{"DISCONNECT,1", Disconnect, "1", nil},
		{"ADDRESS,123,ip1,ipsubnet", Address, "123,ip1,ipsubnet", nil},
		{"UNPARSEABLE", "", "", errors.New("unable to parse event: UNPARSEABLE")},
	}

	for _, test := range testData {
		event, eventData, err := ParseClientEvent(test.testLine)
		assert.Equal(t, test.event, event, test.testLine)
		assert.Equal(t, test.eventData, eventData, test.testLine)
		assert.Equal(t, test.err, err, test.testLine)
	}
}

func TestEnvVarIsParsed(t *testing.T) {
	var testData = []struct {
		testLine string
		key      string
		val      string
		err      error
	}{
		{"abc=123", "abc", "123", nil},
		{"emptyval=", "emptyval", "", nil},
		{"noequal", "noequal", "", nil},
		{"equalsinvalue=abc=123", "equalsinvalue", "abc=123", nil},
	}

	for _, test := range testData {
		key, val, err := ParseEnvVar(test.testLine)
		assert.Equal(t, test.key, key, test.testLine)
		assert.Equal(t, test.val, val, test.testLine)
		assert.Equal(t, test.err, err, test.testLine)
	}
}

func TestIDAndKeyIsParsed(t *testing.T) {
	var testData = []struct {
		testLine string
		ID       int
		key      int
		err      error
	}{
		{"123,456", 123, 456, nil},
		{"abc,def", Undefined, Undefined, errors.New("unable to parse identifiers: abc,def")},
		{"garbage", Undefined, Undefined, errors.New("unable to parse identifiers: garbage")},
		{"123,abc", Undefined, Undefined, errors.New("unable to parse identifiers: 123,abc")},
	}

	for _, test := range testData {
		ID, key, err := ParseIDAndKey(test.testLine)
		assert.Equal(t, test.ID, ID, test.testLine)
		assert.Equal(t, test.key, key, test.testLine)
		assert.Equal(t, test.err, err, test.testLine)
	}
}

func TestIDIsParsed(t *testing.T) {
	var testData = []struct {
		testLine string
		ID       int
		err      error
	}{
		{"123", 123, nil},
		{"garbage", Undefined, errors.New("unable to parse identifier: garbage")},
	}

	for _, test := range testData {
		ID, err := ParseID(test.testLine)
		assert.Equal(t, test.ID, ID, test.testLine)
		assert.Equal(t, test.err, err, test.testLine)
	}

}
