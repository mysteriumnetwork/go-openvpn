package openvpn3

import "sync"

type Logger interface {
	Log(string)
}

type Event struct {
	Fatal bool
	Error bool
	Name  string
	Info  string
}

type EventConsumer interface {
	OnEvent(*Event)
}

type Statistics struct {
	BytesIn  int
	BytesOut int
}

type StatsConsumer interface {
	OnStats(*Statistics)
}

type unregister func()

type callbackRegistry struct {
	lock   sync.Locker
	idMap  map[int]interface{}
	lastId int
}

func (cr *callbackRegistry) register(callbacks interface{}) (int, unregister) {
	cr.lock.Lock()
	defer cr.lock.Unlock()

	id := cr.lastId
	cr.lastId++

	cr.idMap[id] = callbacks

	return id, func() {
		cr.unregister(id)
	}
}

func (cr *callbackRegistry) unregister(id int) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	delete(cr.idMap, id)
}

func (cr *callbackRegistry) log(id int, text string) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	callback, ok := cr.idMap[id]
	if !ok {
		return
	}
	logCallback, ok := callback.(Logger)
	if !ok {
		return
	}
	logCallback.Log(text)
}

func (cr *callbackRegistry) event(id int, event Event) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	callback, ok := cr.idMap[id]
	if !ok {
		return
	}
	eventCallback, ok := callback.(EventConsumer)
	if !ok {
		return
	}
	eventCallback.OnEvent(&event)

}

func (cr *callbackRegistry) stats(id int, stats Statistics) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	callback, ok := cr.idMap[id]
	if !ok {
		return
	}
	statsCallback, ok := callback.(StatsConsumer)
	if !ok {
		return
	}
	statsCallback.OnStats(&stats)

}

func newCallbackRegistry() *callbackRegistry {
	return &callbackRegistry{
		lastId: 0,
		idMap:  make(map[int]interface{}),
		lock:   &sync.Mutex{},
	}
}
