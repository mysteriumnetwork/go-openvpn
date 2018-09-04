package openvpn3

import "testing"
import (
	"github.com/stretchr/testify/assert"
)

func TestTwoCallbacksAreAdded(t *testing.T) {
	registry := NewCallbackRegistry()

	id1, _ := registry.register("abc")
	id2, _ := registry.register("def")

	assert.Equal(t, 0, id1)
	assert.Equal(t, 1, id2)
}

func TestCallbackRemovedOnUnregister(t *testing.T) {
	registry := NewCallbackRegistry()
	id, unregister := registry.register("abc")

	assert.Contains(t, registry.idMap, id)
	unregister()
	assert.NotContains(t, registry.idMap, id)
}

func TestLogCallbackCalled(t *testing.T) {
	registry := NewCallbackRegistry()

	callback := &mockedCallback{}
	id, _ := registry.register(callback)

	registry.Log(id, "test")

	assert.True(t, callback.logCalled)
}

func TestEventCallbackCalled(t *testing.T) {
	registry := NewCallbackRegistry()
	callback := &mockedCallback{}

	id, _ := registry.register(callback)
	registry.Event(id, Event{})

	assert.True(t, callback.eventCalled)
}

func TestStatsCallbackCalled(t *testing.T) {
	registry := NewCallbackRegistry()
	callback := &mockedCallback{}

	id, _ := registry.register(callback)
	registry.Stats(id, Statistics{})

	assert.True(t, callback.statsCalled)
}

type mockedCallback struct {
	logCalled   bool
	eventCalled bool
	statsCalled bool
}

func (mc *mockedCallback) Log(_ string) {
	mc.logCalled = true
}

func (mc *mockedCallback) OnEvent(_ Event) {
	mc.eventCalled = true
}

func (mc *mockedCallback) OnStats(_ Statistics) {
	mc.statsCalled = true
}
