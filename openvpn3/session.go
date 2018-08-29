package openvpn3

/*

#cgo CFLAGS: -I${SRCDIR}/bridge

#cgo LDFLAGS: -L${SRCDIR}/bridge
//main lib link
//TODO reuse GOOS somehow?
#cgo darwin,amd64 LDFLAGS: -lopenvpn3_darwin_amd64
#cgo ios,arm64 LDFLAGS: -lopenvpn3_ios_arm64
#cgo linux LDFLAGS: -lopenvpn3_linux_amd64
#cgo windows LDFLAGS: -lopenvpn3_windows_amd64
//TODO copied from openvpnv3 lib build tool - do we really need all of this?
#cgo darwin,amd64 LDFLAGS: -framework Security -framework CoreFoundation -framework SystemConfiguration -framework IOKit -framework ApplicationServices -mmacosx-version-min=10.8 -stdlib=libc++
//iOS frameworks
#cgo ios,arm64 LDFLAGS: -fobjc-arc -framework UIKit
#cgo windows LDFLAGS: -lws2_32 -liphlpapi

#include <library.h>

extern void GoLogCallback(user_data usrData, char * str);

extern void GoStatsCallback(user_data usrData, conn_stats stats);

extern void GoEventCallback(user_data usrData, conn_event event);
*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

var callbacks = NewCallbackRegistry()

//export GoStatsCallback
func GoStatsCallback(ptr C.user_data, cStats C.conn_stats) {
	id := int(ptr)
	var stats Statistics
	stats.BytesIn = int(cStats.bytes_in)
	stats.BytesOut = int(cStats.bytes_out)
	callbacks.Stats(id, stats)
}

//export GoLogCallback
func GoLogCallback(ptr C.user_data, cStr *C.char) {
	goStr := C.GoString(cStr)
	id := int(ptr)
	callbacks.Log(id, goStr)
}

//export GoEventCallback
func GoEventCallback(ptr C.user_data, cEvent C.conn_event) {
	id := int(ptr)
	var e Event
	e.Error = bool(cEvent.error)
	e.Fatal = bool(cEvent.fatal)
	e.Name = C.GoString(cEvent.name)
	e.Info = C.GoString(cEvent.info)
	callbacks.Event(id, e)
}

type Session struct {
	finished  *sync.WaitGroup
	resError  error
	callbacks interface{}
	session   unsafe.Pointer //handle to created session after Start method is called
}

func SelfCheck(logger Logger) {
	id, callbackRemove := callbacks.Register(logger)
	defer callbackRemove()
	C.check_library(C.user_data(id), C.log_callback(C.GoLogCallback))
	logger.Log("Package version: " + Version())
}

func NewSession(callbacks interface{}) *Session {
	return &Session{
		callbacks: callbacks,
		resError:  nil,
		finished:  &sync.WaitGroup{},
	}
}

type expCredentials C.user_credentials

func (p *Session) Start(profile string, creds Credentials) {
	p.finished.Add(1)
	go func() {
		defer p.finished.Done()

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

		callbackId, removeCallback := callbacks.Register(p.callbacks)
		defer removeCallback()

		session, _ := C.new_session(
			profileContent.Ptr,
			C.user_credentials(cCreds),
			C.user_data(callbackId),
			C.stats_callback(C.GoStatsCallback),
			C.log_callback(C.GoLogCallback),
			C.event_callback(C.GoEventCallback),
		)

		if session == nil {
			p.resError = errors.New("openvpn3 init failed")
			return
		}
		p.session = session

		res, _ := C.start_session(session)
		if res != 0 {
			p.resError = errors.New("openvpn3 connect failed")
		}

		C.cleanup_session(session)
	}()
}

func (p *Session) Wait() error {
	p.finished.Wait()
	return p.resError
}

func (p *Session) Stop() {
	C.stop_session(p.session)
}
