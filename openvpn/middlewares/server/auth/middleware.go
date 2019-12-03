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

package auth

import (
	"strings"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/log"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/middlewares/server"
)

type middleware struct {
	// TODO: consider implementing event channel to communicate required callbacks
	credentialsValidator CredentialsValidator
	commandWriter        management.CommandWriter
	currentEvent         server.ClientEvent
}

// CredentialsValidator callback checks given auth primitives (i.e. customer identity signature / node's sessionId)
type CredentialsValidator func(clientID int, username, password string) (bool, error)

// NewMiddleware creates server user_auth challenge authentication middleware
func NewMiddleware(credentialsValidator CredentialsValidator) *middleware {
	return &middleware{
		credentialsValidator: credentialsValidator,
		commandWriter:        nil,
		currentEvent:         server.UndefinedEvent,
	}
}

func (m *middleware) Start(commandWriter management.CommandWriter) error {
	m.commandWriter = commandWriter
	return nil
}

func (m *middleware) Stop(commandWriter management.CommandWriter) error {
	return nil
}

func (m *middleware) ConsumeLine(line string) (bool, error) {
	if !strings.HasPrefix(line, ">CLIENT:") {
		return false, nil
	}

	clientLine := strings.TrimPrefix(line, ">CLIENT:")

	eventType, eventData, err := server.ParseClientEvent(clientLine)
	if err != nil {
		return true, err
	}

	switch eventType {
	case server.Connect, server.Reauth:
		ID, key, err := server.ParseIDAndKey(eventData)
		if err != nil {
			return true, err
		}
		m.startOfEvent(eventType, ID, key)
	case server.Env:
		if strings.ToLower(eventData) == "end" {
			m.endOfEvent()
			return true, nil
		}

		key, val, err := server.ParseEnvVar(eventData)
		if err != nil {
			return true, err
		}
		m.addEnvVar(key, val)
	case server.Established, server.Disconnect:
		ID, err := server.ParseID(eventData)
		if err != nil {
			return true, err
		}
		m.startOfEvent(eventType, ID, server.Undefined)
	case server.Address:
		log.Info("Address for client:", eventData)
	default:
		log.Error("Undefined user notification event:", eventType, eventData)
		log.Error("Original line was:", line)
	}
	return true, nil
}

func (m *middleware) startOfEvent(eventType server.ClientEventType, clientID int, keyID int) {
	m.currentEvent.EventType = eventType
	m.currentEvent.ClientID = clientID
	m.currentEvent.ClientKey = keyID
}

func (m *middleware) addEnvVar(key string, val string) {
	m.currentEvent.Env[key] = val
}

func (m *middleware) endOfEvent() {
	m.handleClientEvent(m.currentEvent)
	m.reset()
}

func (m *middleware) reset() {
	m.currentEvent = server.UndefinedEvent
}

func (m *middleware) handleClientEvent(event server.ClientEvent) {
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
		// NOTE: do not cleanup session after disconnect event risen by transport itself
		//  cleanup session only by user's intent
	}
}

func (m *middleware) authenticateClient(clientID, clientKey int, username, password string) error {

	if username == "" || password == "" {
		return denyClientAuthWithMessage(m.commandWriter, clientID, clientKey, "missing username or password")
	}

	log.Info("Authenticating user:", username, "clientID:", clientID, "clientKey:", clientKey)

	authenticated, err := m.credentialsValidator(clientID, username, password)
	if err != nil {
		log.Error("Authentication error:", err)
		return denyClientAuthWithMessage(m.commandWriter, clientID, clientKey, "internal error")
	}

	if authenticated {
		return approveClient(m.commandWriter, clientID, clientKey)
	}
	return denyClientAuthWithMessage(m.commandWriter, clientID, clientKey, "wrong username or password")
}

func approveClient(commandWriter management.CommandWriter, clientID, keyID int) error {
	_, err := commandWriter.SingleLineCommand("client-auth-nt %d %d", clientID, keyID)
	return err
}

func denyClientAuthWithMessage(commandWriter management.CommandWriter, clientID, keyID int, message string) error {
	_, err := commandWriter.SingleLineCommand("client-deny %d %d %s", clientID, keyID, message)
	return err
}
