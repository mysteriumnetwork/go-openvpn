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
)

// Primitives structure holds TLS primitives required to setup basic cryptographics for openvpn server/client
type Primitives struct {
	CertificateAuthority *CertificateAuthority
	ServerCertificate    *CertificateKeyPair
	PresharedKey         *TLSPresharedKey
}

// NewTLSPrimitives function creates TLS primitives for given service location and provider id
func NewTLSPrimitives(caCertSubject, serverCertSubject pkix.Name) (*Primitives, error) {

	key, err := createTLSCryptKey()
	if err != nil {
		return nil, err
	}

	ca, err := CreateAuthority(newCACert(caCertSubject))
	if err != nil {
		return nil, err
	}

	server, err := ca.CreateDerived(newServerCert(x509.ExtKeyUsageServerAuth, serverCertSubject))
	if err != nil {
		return nil, err
	}

	return &Primitives{
		CertificateAuthority: ca,
		ServerCertificate:    server,
		PresharedKey:         &key,
	}, nil
}
