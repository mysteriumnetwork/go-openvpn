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

package openvpn3

import "sync"

type tunSetupRegistry struct {
	lock   sync.Locker
	idMap  map[int]TunnelSetup
	lastId int
}

func (registry *tunSetupRegistry) register(delegate TunnelSetup) (int, func()) {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	id := registry.lastId
	registry.lastId++
	registry.idMap[id] = delegate

	return id, func() {
		registry.unregister(id)
	}
}

func (registry *tunSetupRegistry) unregister(id int) {
	registry.lock.Lock()
	defer registry.lock.Unlock()
	delete(registry.idMap, id)
}

func (registry *tunSetupRegistry) lookup(id int) TunnelSetup {
	//to avoid error handling in Go -> C callbacks, return Noop TunnelSetup in case id is not registered or already
	//removed
	registry.lock.Lock()
	defer registry.lock.Unlock()
	tunSetup, ok := registry.idMap[id]
	if !ok {
		return &NoOpTunnelSetup{}
	}
	return tunSetup
}
