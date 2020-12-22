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

package tls

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// TLSPresharedKey defines TLS preshared key type
type TLSPresharedKey []byte

// ToPEMFormat renders TLS preshared key to PEM compatible string which can be written to PEM file
func (key TLSPresharedKey) ToPEMFormat() string {
	buffer := bytes.Buffer{}

	fmt.Fprintln(&buffer, "-----BEGIN OpenVPN Static key V1-----")
	fmt.Fprintln(&buffer, hex.EncodeToString(key))
	fmt.Fprintln(&buffer, "-----END OpenVPN Static key V1-----")

	return buffer.String()
}

// createTLSCryptKey generates symmetric key in HEX format 2048 bits length
func createTLSCryptKey() (TLSPresharedKey, error) {

	taKey := make([]byte, 256)
	_, err := rand.Read(taKey)
	if err != nil {
		return nil, err
	}
	return TLSPresharedKey(taKey), nil
}
