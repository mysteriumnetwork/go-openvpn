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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)

// CertificateKeyPair represents x509 type certificate and corresponding private key
type CertificateKeyPair struct {
	privateKey *ecdsa.PrivateKey
	x509cert   *x509.Certificate
	certBytes  []byte
	keyBytes   []byte
}

// ToPEMFormat method returns certificate serialized to string by PEM encoding rules
func (ckp *CertificateKeyPair) ToPEMFormat() string {
	return string(pem.EncodeToMemory(pemBlock("CERTIFICATE", ckp.certBytes)))
}

// KeyToPEMFormat returns private key serialized to string by PEM encoding rules
func (ckp *CertificateKeyPair) KeyToPEMFormat() string {
	return string(pem.EncodeToMemory(pemBlock("EC PRIVATE KEY", ckp.keyBytes)))
}

func pemBlock(blockType string, data []byte) *pem.Block {
	return &pem.Block{
		Type:  blockType,
		Bytes: data,
	}
}

// CertificateAuthority represents self-signed certificate/key pair which can create signed derived certificates
type CertificateAuthority struct {
	CertificateKeyPair
}

// CreateDerived creates new certificate/key by given x509 data and signed by current authority
func (ca *CertificateAuthority) CreateDerived(template *x509.Certificate) (*CertificateKeyPair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, ca.x509cert, &privateKey.PublicKey, ca.privateKey)
	if err != nil {
		return nil, err
	}
	return &CertificateKeyPair{
		privateKey: privateKey,
		x509cert:   template,
		certBytes:  certBytes,
		keyBytes:   keyBytes,
	}, nil
}

// CreateAuthority creates new self signed certificate with given x509 data
func CreateAuthority(template *x509.Certificate) (*CertificateAuthority, error) {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	return &CertificateAuthority{
		CertificateKeyPair{
			privateKey: privateKey,
			x509cert:   template,
			certBytes:  certBytes,
			keyBytes:   keyBytes,
		},
	}, nil
}
