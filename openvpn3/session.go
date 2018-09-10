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

type Session struct {
	finished   *sync.WaitGroup
	resError   error
	callbacks  interface{}
	sessionPtr unsafe.Pointer //handle to created sessionPtr after Start method is called
}

func NewSession(callbacks interface{}) *Session {
	return &Session{
		callbacks: callbacks,
		resError:  nil,
		finished:  &sync.WaitGroup{},
	}
}

var ErrInitFailed = errors.New("openvpn3 init failed")
var ErrConnectFailed = errors.New("openvpn3 connect failed")

type expCredentials C.user_credentials

func (session *Session) Start(profile string, creds Credentials) {
	session.finished.Add(1)
	go func() {
		defer session.finished.Done()

		profileContent := NewCharPointer(profile)
		defer profileContent.Delete()

		cUsername := NewCharPointer(creds.Username)
		defer cUsername.Delete()

		cPassword := NewCharPointer(creds.Password)
		defer cPassword.Delete()

		cCreds := expCredentials{
			username: cUsername.Ptr,
			password: cPassword.Ptr,
		}

		callbacksDelegate, removeCallback := registerCallbackDelegate(session.callbacks)
		defer removeCallback()

		tunBuilderCallbacks, removeTunCallbacks := registerTunnelSetupDelegate(&NoOpTunnelSetup{})
		defer removeTunCallbacks()

		sessionPtr, _ := C.new_session(
			profileContent.Ptr,
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

func (session *Session) Wait() error {
	session.finished.Wait()
	return session.resError
}

func (session *Session) Stop() {
	C.stop_session(session.sessionPtr)
}
