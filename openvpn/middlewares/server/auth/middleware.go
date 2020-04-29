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
	"sync"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/log"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/middlewares/server"
)

// ClientEventCallback is called when state of each OpenVPN client changes.
type ClientEventCallback func(event server.ClientEvent)

// Middleware is able to subscribe to client status events, exposes client control API.
//
// The OpenVPN server should have been started with the
// --management-client-auth directive so that it will ask the management
// interface to approve client connections.
type Middleware struct {
	commandWriter management.CommandWriter
	currentEvent  server.ClientEvent

	listenersMu sync.RWMutex
	listeners   []ClientEventCallback
}

// NewMiddleware creates new instance of Middleware.
func NewMiddleware(listeners ...ClientEventCallback) *Middleware {
	return &Middleware{
		currentEvent: server.UndefinedEvent(),
		listeners:    listeners,
	}
}

// ClientsSubscribe subscribes to Openvpn clients states.
func (m *Middleware) ClientsSubscribe(callback ClientEventCallback) {
	m.listenersMu.Lock()
	defer m.listenersMu.Unlock()

	m.listeners = append(m.listeners, callback)
}

// ClientAccept is a client control which allows authorization (for CONNECT or REAUTH state).
func (m *Middleware) ClientAccept(clientID, keyID int) error {
	_, err := m.commandWriter.SingleLineCommand("client-auth-nt %d %d", clientID, keyID)
	return err
}

// ClientDeny is a client control which forbids authorization (for CONNECT or REAUTH state).
func (m *Middleware) ClientDeny(clientID, keyID int, message string) error {
	_, err := m.commandWriter.SingleLineCommand("client-deny %d %d", clientID, keyID, message)
	return err
}

// ClientDenyWithMessage is a client control which forbids authorization with reason message (for CONNECT or REAUTH state).
func (m *Middleware) ClientDenyWithMessage(clientID, keyID int, message string) error {
	_, err := m.commandWriter.SingleLineCommand("client-deny %d %d %s", clientID, keyID, message)
	return err
}

// ClientKill is a client control which stops established connection (for ESTABLISHED state).
func (m *Middleware) ClientKill(clientID int) error {
	_, err := m.commandWriter.SingleLineCommand("client-kill %d", clientID)
	return err
}

// ClientKillWithMessage is a client control which stops established connection with reason message (for ESTABLISHED state).
func (m *Middleware) ClientKillWithMessage(clientID int, message string) error {
	_, err := m.commandWriter.SingleLineCommand("client-kill %d %s", clientID, message)
	return err
}

// Start starts the middleware.
func (m *Middleware) Start(commandWriter management.CommandWriter) error {
	m.commandWriter = commandWriter
	return nil
}

// Stop stops the middleware.
func (m *Middleware) Stop(_ management.CommandWriter) error {
	return nil
}

// ConsumeLine handles the given openvpn management line.
func (m *Middleware) ConsumeLine(line string) (bool, error) {
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

func (m *Middleware) startOfEvent(eventType server.ClientEventType, clientID int, keyID int) {
	m.currentEvent.EventType = eventType
	m.currentEvent.ClientID = clientID
	m.currentEvent.ClientKey = keyID
}

func (m *Middleware) addEnvVar(key string, val string) {
	m.currentEvent.Env[key] = val
}

func (m *Middleware) endOfEvent() {
	m.listenersMu.RLock()
	defer m.listenersMu.RUnlock()

	for _, subscription := range m.listeners {
		subscription(m.currentEvent)
	}
	m.reset()
}

func (m *Middleware) reset() {
	m.currentEvent = server.UndefinedEvent()
}
