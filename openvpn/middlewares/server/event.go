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

package server

// ClientEventType is a OpenVPN management client event.
type ClientEventType string

const (
	// Connect represent a CONNECT OpenVPN event type.
	Connect = ClientEventType("CONNECT")
	// Reauth represent a REAUTH OpenVPN event type.
	Reauth = ClientEventType("REAUTH")
	// Established represent a ESTABLISHED OpenVPN event type.
	Established = ClientEventType("ESTABLISHED")
	// Disconnect represent a DISCONNECT OpenVPN event type.
	Disconnect = ClientEventType("DISCONNECT")
	// Address represent a ADDRESS OpenVPN event type.
	Address = ClientEventType("ADDRESS")

	//Env is a pseudo event type ENV - that means some of above defined events are multiline and ENV messages are part of it
	Env = ClientEventType("ENV")
	//Undefined is a constant which means that id of type int is undefined
	Undefined = -1
)

// ClientEvent represent a OpenVPN management client event.
type ClientEvent struct {
	EventType ClientEventType
	ClientID  int
	ClientKey int
	Env       map[string]string
}

// UndefinedEvent is an empty OpenVPN management client event.
// Note: can't be a var, because `ClientEvent.Env` variable is being overwritten as it's a reference to same memory.
func UndefinedEvent() ClientEvent {
	return ClientEvent{
		ClientID:  Undefined,
		ClientKey: Undefined,
		Env:       make(map[string]string),
	}
}
