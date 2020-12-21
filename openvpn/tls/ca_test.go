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
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	"github.com/stretchr/testify/assert"
)

var caSubject = pkix.Name{
	Country:            []string{"GB"},
	Organization:       []string{""},
	OrganizationalUnit: []string{""},
}

var serverCertSubject = pkix.Name{
	Country:            []string{"GB"},
	Organization:       []string{""},
	OrganizationalUnit: []string{""},
	CommonName:         "some fake identity ",
}

func TestCertificateAuthorityIsCreatedAndCertCanBeSerialized(t *testing.T) {
	_, err := CreateAuthority(newCACert(caSubject))
	assert.NoError(t, err)
}

func TestServerCertificateIsCreatedAndBothCertAndKeyCanBeSerialized(t *testing.T) {
	ca, err := CreateAuthority(newCACert(caSubject))
	assert.NoError(t, err)
	_, err = ca.CreateDerived(newServerCert(x509.ExtKeyUsageServerAuth, serverCertSubject))
	assert.NoError(t, err)
}
