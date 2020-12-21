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

package tls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedKey = `-----BEGIN OpenVPN Static key V1-----
616263
-----END OpenVPN Static key V1-----
`

func TestTLSPresharedKeyProducesValidPEMFormat(t *testing.T) {
	key := TLSPresharedKey("abc")
	assert.Equal(
		t,
		expectedKey,
		key.ToPEMFormat(),
	)
}

func TestGeneratedKeyIsExpectedSize(t *testing.T) {
	key, err := createTLSCryptKey()
	assert.NoError(t, err)
	assert.Len(t, key, 256)
}
