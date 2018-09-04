package openvpn3

import "sync"

type tunSetupRegistry struct {
	lock   sync.Locker
	idMap  map[int]TunnelSetup
	lastId int
}

func (registry *tunSetupRegistry) Register(delegate TunnelSetup) (int, func()) {
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
