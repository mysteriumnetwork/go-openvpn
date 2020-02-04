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

package bytecount

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	sbc SessionByteCount
}

func (mh *mockHandler) Handle(s SessionByteCount) {
	mh.sbc = s
}

func Test_ConsumesOKLine(t *testing.T) {
	statsRecorder := &mockHandler{}
	middleware := NewMiddleware(statsRecorder.Handle, 1)
	consumed, err := middleware.ConsumeLine(">BYTECOUNT_CLI:1,2,3")
	assert.Nil(t, err)
	assert.True(t, consumed)
	assert.Equal(t, 1, statsRecorder.sbc.ClientID)
	assert.Equal(t, uint64(2), statsRecorder.sbc.BytesIn)
	assert.Equal(t, uint64(3), statsRecorder.sbc.BytesOut)
}

func Test_IgnoresMalformedLines(t *testing.T) {
	badLines := []string{
		">BYTECOUNT_CLI:asdasd,2,3",
		">BYTECOUNT_CLI",
		"whatever",
		"BYTECOUNT_CLI:1,2,3",
	}
	for _, v := range badLines {
		statsRecorder := &mockHandler{}
		middleware := NewMiddleware(statsRecorder.Handle, 1)
		consumed, err := middleware.ConsumeLine(v)
		assert.Nil(t, err)
		assert.False(t, consumed)
		assert.EqualValues(t, SessionByteCount{}, statsRecorder.sbc)
	}
}
