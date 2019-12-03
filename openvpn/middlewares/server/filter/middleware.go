/*
 * Copyright (C) 2019 The "MysteriumNetwork/go-openvpn" Authors.
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

package filter

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/log"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/middlewares/server"
)

const filterLANTemplate = `client-pf {{.ClientID}}
[CLIENTS DROP]
[SUBNETS ACCEPT]
{{- range $subnet := .Allow}}
+{{$subnet}}
{{- end}}
{{- range $subnet := .Block}}
-{{$subnet}}
{{- end}}
[END]
END
`

var filterLAN = template.Must(template.New("filter_lan").Parse(filterLANTemplate))

type middleware struct {
	commandWriter management.CommandWriter
	currentEvent  server.ClientEvent
	allow         []string
	block         []string
}

// NewMiddleware creates server user_auth challenge authentication middleware
func NewMiddleware(allow, block []string) *middleware {
	return &middleware{
		commandWriter: nil,
		allow:         allow,
		block:         block,
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
		clientID, _, err := server.ParseIDAndKey(eventData)
		if err != nil {
			return true, err
		}

		m.currentEvent.EventType = eventType
		m.currentEvent.ClientID = clientID
	case server.Env:
		if strings.ToLower(eventData) == "end" {
			m.handleClientEvent(m.currentEvent)
		}
	}

	return true, nil
}

func (m *middleware) handleClientEvent(event server.ClientEvent) {
	switch event.EventType {
	case server.Connect, server.Reauth:
		if err := filterSubnets(m.commandWriter, event.ClientID, m.allow, m.block); err != nil {
			log.Error("Unable to authenticate client:", err)
		}
	}
}

func filterSubnets(commandWriter management.CommandWriter, clientID int, allow, block []string) error {
	data := struct {
		ClientID int
		Allow    []string
		Block    []string
	}{
		ClientID: clientID,
		Allow:    allow,
		Block:    block,
	}

	var tpl bytes.Buffer
	if err := filterLAN.Execute(&tpl, data); err != nil {
		return err
	}

	_, err := commandWriter.SingleLineCommand(tpl.String())

	return err
}
