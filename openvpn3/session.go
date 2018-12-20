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

/*

#cgo CFLAGS: -I${SRCDIR}/bridge

#cgo LDFLAGS: -L${SRCDIR}/bridge
//main lib link
//TODO reuse GOOS somehow?
#cgo darwin,amd64 LDFLAGS: -lopenvpn3_darwin_amd64
#cgo ios,arm64 LDFLAGS: -lopenvpn3_ios_arm64
#cgo linux,!android,amd64 LDFLAGS: -lopenvpn3_linux_amd64
#cgo windows LDFLAGS: -lopenvpn3_windows_amd64
#cgo android,arm64 LDFLAGS: -lopenvpn3_android_arm64
#cgo android,amd64 LDFLAGS: -lopenvpn3_android_amd64
#cgo android,386 LDFLAGS: -lopenvpn3_android_x86
#cgo android,arm LDFLAGS: -lopenvpn3_android_armeabi-v7a
//TODO copied from openvpnv3 lib build tool - do we really need all of this?
#cgo darwin,amd64 LDFLAGS: -framework Security -framework CoreFoundation -framework SystemConfiguration -framework IOKit -framework ApplicationServices -mmacosx-version-min=10.8 -stdlib=libc++
//iOS frameworks
#cgo ios,arm64 LDFLAGS: -fobjc-arc -framework UIKit
#cgo windows LDFLAGS: -lws2_32 -liphlpapi

#include <library.h>

*/
import "C"
import (
	"errors"
	"sync"
)

// Session represents the openvpn session
type Session struct {
	config          Config
	userCredentials UserCredentials
	callbacks       interface{}
	tunnelSetup     TunnelSetup
	finished        sync.WaitGroup
	stop            sync.WaitGroup
	reconnectChan   chan int

	// runtime variables
	resError error
}

// NewSession creates a new session given the callbacks
func NewSession(config Config, userCredentials UserCredentials, callbacks interface{}) *Session {
	return &Session{
		config:          config,
		userCredentials: userCredentials,
		callbacks:       callbacks,
		tunnelSetup:     &NoOpTunnelSetup{},
		resError:        nil,
		reconnectChan:   make(chan int, 1),
	}
}

// MobileSessionCallbacks are the callbacks required for a mobile session
type MobileSessionCallbacks interface {
	EventConsumer
	Logger
	StatsConsumer
}

// NewMobileSession creates a new mobile session provided the required callbacks and tunnel setup
func NewMobileSession(config Config, userCredentials UserCredentials, callbacks MobileSessionCallbacks, tunSetup TunnelSetup) *Session {
	return &Session{
		config:          config,
		userCredentials: userCredentials,
		callbacks:       callbacks,
		tunnelSetup:     tunSetup,
		resError:        nil,
		reconnectChan:   make(chan int, 1),
	}
}

// ErrInitFailed is the error we return when openvpn3 initiation fails
var ErrInitFailed = errors.New("openvpn3 init failed")

// ErrConnectFailed is the error we return when openvpn3 fails to connect
var ErrConnectFailed = errors.New("openvpn3 connect failed")

// Start starts the session
func (session *Session) Start() {
	session.finished.Add(1)
	session.stop.Add(1)
	go func() {
		defer session.finished.Done()
		cConfig, cConfigUnregister := session.config.toPtr()
		defer cConfigUnregister()

		cCredentials, cCredentialsUnregister := session.userCredentials.toPtr()
		defer cCredentialsUnregister()

		callbacksDelegate, removeCallback := registerCallbackDelegate(session.callbacks)
		defer removeCallback()

		tunBuilderCallbacks, removeTunCallbacks := registerTunnelSetupDelegate(session.tunnelSetup)
		defer removeTunCallbacks()

		sessionPtr, _ := C.new_session(
			cConfig,
			cCredentials,
			C.callbacks_delegate(callbacksDelegate),
			C.tun_builder_callbacks(tunBuilderCallbacks),
		)

		if sessionPtr == nil {
			session.resError = ErrInitFailed
			return
		}

		go func() {
			session.stop.Wait()
			C.stop_session(sessionPtr)
		}()

		go func() {
			for seconds := range session.reconnectChan {
				C.reconnect_session(sessionPtr, C.int(seconds))
			}
		}()

		res, _ := C.start_session(sessionPtr)
		if res != 0 {
			session.resError = ErrConnectFailed
		}

		C.cleanup_session(sessionPtr)
	}()
}

// Wait waits for the session to complete
func (session *Session) Wait() error {
	session.finished.Wait()
	return session.resError
}

// Stop stops the session
func (session *Session) Stop() {
	session.stop.Done()
	close(session.reconnectChan)
}

// Reconnect session without propagating DISCONNECT event after `seconds` time
func (session *Session) Reconnect(seconds int) error {
	session.reconnectChan <- seconds
	return nil
}
