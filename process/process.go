package process

/*

#cgo CFLAGS: -I${SRCDIR}/bridge
#cgo LDFLAGS: -lstdc++
#cgo LDFLAGS: -L${SRCDIR}/bridge -lopenvpn
//TODO copied from openvpnv3 lib build tool - do we really need all of this?
#cgo darwin LDFLAGS: -framework Security -framework CoreFoundation -framework SystemConfiguration -framework IOKit -framework ApplicationServices

#include <process.h>

extern void GoLogCallback(user_data usrData, char * str);

extern void GoStatsCallback(user_data usrData, conn_stats stats);

extern void GoEventCallback(user_data usrData, conn_event event);
*/
import "C"
import (
	"fmt"
	"sync"
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

type Process struct {
	finished  *sync.WaitGroup
	resError  error
	callbacks interface{}
}

func CheckLibrary(logger Logger) {
	id, callbackRemove := callbacks.Register(logger)
	defer callbackRemove()
	C.checkLibrary(C.user_data(id), C.log_callback(C.GoLogCallback))
}

func NewProcess(callbacks interface{}) *Process {
	return &Process{
		callbacks: callbacks,
		resError:  nil,
		finished:  &sync.WaitGroup{},
	}
}

func (p *Process) RunWithArgs(args ...string) {
	p.finished.Add(1)
	go func() {
		defer p.finished.Done()

		cPtr := NewCharPointer(args[0])
		defer cPtr.Delete()

		callbackId, removeCallback := callbacks.Register(p.callbacks)
		defer removeCallback()

		res, err := C.initProcess(
			cPtr.Ptr,
			C.user_data(callbackId),
			C.stats_callback(C.GoStatsCallback),
			C.log_callback(C.GoLogCallback),
			C.event_callback(C.GoEventCallback),
		)
		if err != nil {
			p.resError = err
		} else if res != 0 {
			p.resError = fmt.Errorf("res error: %v", res)
		} else {
			p.resError = nil
		}
	}()
}

func (p *Process) WaitFor() error {
	p.finished.Wait()
	return p.resError
}

func (p *Process) Stop() {

}
