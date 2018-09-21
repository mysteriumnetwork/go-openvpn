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
#cgo linux,amd64 LDFLAGS: -lopenvpn3_linux_amd64
#cgo windows LDFLAGS: -lopenvpn3_windows_amd64
#cgo android,arm64 LDFLAGS: -lopenvpn3_android_arm64
//TODO copied from openvpnv3 lib build tool - do we really need all of this?
#cgo darwin,amd64 LDFLAGS: -framework Security -framework CoreFoundation -framework SystemConfiguration -framework IOKit -framework ApplicationServices -mmacosx-version-min=10.8 -stdlib=libc++
//iOS frameworks
#cgo ios,arm64 LDFLAGS: -fobjc-arc -framework UIKit
#cgo windows LDFLAGS: -lws2_32 -liphlpapi

#include <library.h>
#include <tunsetup.h>

*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

// Session represents the openvpn session
type Session struct {
	finished    *sync.WaitGroup
	resError    error
	callbacks   interface{}
	tunnelSetup TunnelSetup
	sessionPtr  unsafe.Pointer //handle to created sessionPtr after Start method is called
}

// NewSession creates a new session given the callbacks
func NewSession(callbacks interface{}) *Session {
	return &Session{
		callbacks:   callbacks,
		resError:    nil,
		finished:    &sync.WaitGroup{},
		tunnelSetup: &NoOpTunnelSetup{},
	}
}

// MobileSessionCallbacks are the callbacks required for a mobile session
type MobileSessionCallbacks interface {
	EventConsumer
	Logger
	StatsConsumer
}

// NewMobileSession creates a new mobile session provided the required callbacks and tunnel setup
func NewMobileSession(callbacks MobileSessionCallbacks, tunSetup TunnelSetup) *Session {
	return &Session{
		callbacks:   callbacks,
		resError:    nil,
		finished:    &sync.WaitGroup{},
		tunnelSetup: tunSetup,
	}
}

// ErrInitFailed is the error we return when openvpn3 initiation fails
var ErrInitFailed = errors.New("openvpn3 init failed")

// ErrConnectFailed is the error we return when openvpn3 fails to connect
var ErrConnectFailed = errors.New("openvpn3 connect failed")

type expConfig C.config
type expUserCredentials C.user_credentials

// Start starts the session
func (session *Session) Start(profile string, creds Credentials) {
	session.finished.Add(1)
	go func() {
		defer session.finished.Done()

		profileContent := newCharPointer(profile)
		defer profileContent.delete()

		guiVersion := newCharPointer("cli 1.0")
		defer guiVersion.delete()

		compressionMode := newCharPointer("yes")
		defer compressionMode.delete()

		cConfig := expConfig{
			profileContent:    profileContent.Ptr,
			guiVersion:        guiVersion.Ptr,
			info:              true,
			clockTickMS:       1000, // ticks every 1 sec
			disableClientCert: true,
			connTimeout:       10, // 10 seconds
			tunPersist:        true,
			compressionMode:   compressionMode.Ptr,
		}

		cUsername := newCharPointer(creds.Username)
		defer cUsername.delete()

		cPassword := newCharPointer(creds.Password)
		defer cPassword.delete()

		cCreds := expUserCredentials{
			username: cUsername.Ptr,
			password: cPassword.Ptr,
		}

		callbacksDelegate, removeCallback := registerCallbackDelegate(session.callbacks)
		defer removeCallback()

		tunBuilderCallbacks, removeTunCallbacks := registerTunnelSetupDelegate(&NoOpTunnelSetup{})
		defer removeTunCallbacks()

		sessionPtr, _ := C.new_session(
			C.config(cConfig),
			C.user_credentials(cCreds),
			C.callbacks_delegate(callbacksDelegate),
			C.tun_builder_callbacks(tunBuilderCallbacks),
		)

		if sessionPtr == nil {
			session.resError = ErrInitFailed
			return
		}
		session.sessionPtr = sessionPtr

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
	C.stop_session(session.sessionPtr)
}
