/*
 * Copyright (C) 2020 The "MysteriumNetwork/go-openvpn" Authors.
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

package credentials

import (
	"github.com/mysteriumnetwork/go-openvpn/openvpn/log"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/middlewares/server"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/middlewares/server/auth"
)

// Middleware is able to authorize incoming clients by given credentials validator callback.
type Middleware struct {
	*auth.Middleware

	validator Validator
}

// Validator callback checks given auth primitives.
type Validator func(clientID int, username, password string) (bool, error)

// NewMiddleware creates server user_auth challenge authentication Middleware
func NewMiddleware(validator Validator) *Middleware {
	m := new(Middleware)
	m.Middleware = auth.NewMiddleware(m.handleClientEvent)
	m.validator = validator
	return m
}

func (m *Middleware) handleClientEvent(event server.ClientEvent) {
	switch event.EventType {
	case server.Connect, server.Reauth:
		username := event.Env["username"]
		password := event.Env["password"]
		err := m.authenticateClient(event.ClientID, event.ClientKey, username, password)
		if err != nil {
			log.Error("Unable to authenticate client:", err)
		}
	case server.Established:
		log.Info("Client with ID:", event.ClientID, "connection established successfully")
	case server.Disconnect:
		log.Info("Client with ID:", event.ClientID, "disconnected")
	}
}

func (m *Middleware) authenticateClient(clientID, clientKey int, username, password string) error {
	log.Info("Authenticating user:", username, "clientID:", clientID, "clientKey:", clientKey)
	if username == "" || password == "" {
		return m.ClientDenyWithMessage(clientID, clientKey, "missing username or password")
	}

	authenticated, err := m.validator(clientID, username, password)
	if err != nil {
		log.Error("Authentication error:", err)
		return m.ClientDenyWithMessage(clientID, clientKey, "internal error")
	}

	if !authenticated {
		return m.ClientDenyWithMessage(clientID, clientKey, "wrong username or password")
	}

	return m.ClientAccept(clientID, clientKey)
}
