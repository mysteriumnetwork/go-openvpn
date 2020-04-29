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

	"github.com/mysteriumnetwork/go-openvpn/openvpn/log"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/management"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/middlewares/server"
	"github.com/mysteriumnetwork/go-openvpn/openvpn/middlewares/server/auth"
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

// Exposes API to control client's packet filtering.
//
// The OpenVPN server should have been started with the
// --management-client-pf directive so that it will require that
// VPN tunnel packets sent or received by client instances must
// conform to that client's packet filter configuration.
type middleware struct {
	*auth.Middleware

	commandWriter management.CommandWriter
	allow         []string
	block         []string
}

// NewMiddleware creates new instance of middleware
func NewMiddleware(allow, block []string) *middleware {
	m := new(middleware)
	m.Middleware = auth.NewMiddleware(m.handleClientEvent)
	m.allow = allow
	m.block = block
	return m
}

func (m *middleware) Start(commandWriter management.CommandWriter) error {
	m.commandWriter = commandWriter
	return m.Middleware.Start(commandWriter)
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
