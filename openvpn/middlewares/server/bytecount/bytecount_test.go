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
