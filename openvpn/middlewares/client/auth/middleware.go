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
	"regexp"

	"github.com/mysteriumnetwork/go-openvpn/openvpn"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/log"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
)

// CredentialsProvider returns client's current auth primitives (i.e. customer identity signature / node's sessionId)
type CredentialsProvider func() (username string, password string, err error)

type middleware struct {
	fetchCredentials CredentialsProvider
	commandWriter    management.CommandWriter
	lastUsername     string
	lastPassword     string
	state            openvpn.State
}

var rule = regexp.MustCompile("^>PASSWORD:Need 'Auth' username/password$")

// NewMiddleware creates client user_auth challenge authentication middleware
func NewMiddleware(credentials CredentialsProvider) *middleware {
	return &middleware{
		fetchCredentials: credentials,
		commandWriter:    nil,
	}
}

func (m *middleware) Start(commandWriter management.CommandWriter) error {
	m.commandWriter = commandWriter
	log.Info("Starting client user-pass provider middleware")
	return nil
}

func (m *middleware) Stop(connection management.CommandWriter) error {
	return nil
}

func (m *middleware) ConsumeLine(line string) (consumed bool, err error) {
	match := rule.FindStringSubmatch(line)
	if len(match) == 0 {
		return false, nil
	}

	username, password, err := m.fetchCredentials()
	if err != nil {
		return false, err
	}

	log.Info("Authenticating user", username)

	_, err = m.commandWriter.SingleLineCommand("password 'Auth' %s", password)
	if err != nil {
		return true, err
	}

	_, err = m.commandWriter.SingleLineCommand("username 'Auth' %s", username)
	if err != nil {
		return true, err
	}
	return true, nil
}
