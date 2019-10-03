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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwoCallbacksAreAdded(t *testing.T) {
	registry := newCallbackRegistry()

	id1, _ := registry.register("abc")
	id2, _ := registry.register("def")

	assert.Equal(t, 0, id1)
	assert.Equal(t, 1, id2)
}

func TestCallbackRemovedOnUnregister(t *testing.T) {
	registry := newCallbackRegistry()
	id, unregister := registry.register("abc")

	assert.Contains(t, registry.idMap, id)
	unregister()
	assert.NotContains(t, registry.idMap, id)
}

func TestLogCallbackCalled(t *testing.T) {
	registry := newCallbackRegistry()

	callback := &mockedCallback{}
	id, _ := registry.register(callback)

	registry.log(id, "test")

	assert.True(t, callback.logCalled)
}

func TestEventCallbackCalled(t *testing.T) {
	registry := newCallbackRegistry()
	callback := &mockedCallback{}

	id, _ := registry.register(callback)
	registry.event(id, Event{})

	assert.True(t, callback.eventCalled)
}

func TestStatsCallbackCalled(t *testing.T) {
	registry := newCallbackRegistry()
	callback := &mockedCallback{}

	id, _ := registry.register(callback)
	registry.stats(id, Statistics{})

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
