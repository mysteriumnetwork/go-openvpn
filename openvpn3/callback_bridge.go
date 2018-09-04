package openvpn3

/*

#cgo CFLAGS: -I${SRCDIR}/bridge

#include <library.h>
#include <tunsetup.h>

extern void goLogCallback(user_callback_data usrData, char * str);

extern void goStatsCallback(user_callback_data usrData, conn_stats stats);

extern void goEventCallback(user_callback_data usrData, conn_event event);
*/
import "C"

var callbacks = NewCallbackRegistry()

//export goStatsCallback
func goStatsCallback(ptr C.user_callback_data, cStats C.conn_stats) {
	id := int(ptr)
	var stats Statistics
	stats.BytesIn = int(cStats.bytes_in)
	stats.BytesOut = int(cStats.bytes_out)
	callbacks.Stats(id, stats)
}

//export goLogCallback
func goLogCallback(ptr C.user_callback_data, cStr *C.char) {
	goStr := C.GoString(cStr)
	id := int(ptr)
	callbacks.Log(id, goStr)
}

//export goEventCallback
func goEventCallback(ptr C.user_callback_data, cEvent C.conn_event) {
	id := int(ptr)
	var e Event
	e.Error = bool(cEvent.error)
	e.Fatal = bool(cEvent.fatal)
	e.Name = C.GoString(cEvent.name)
	e.Info = C.GoString(cEvent.info)
	callbacks.Event(id, e)
}

func SelfCheck(logger Logger) {
	id, callbackRemove := callbacks.register(logger)
	defer callbackRemove()
	C.check_library(C.user_callback_data(id), C.log_callback(C.goLogCallback))
}

type expCallbacks C.callbacks_delegate

func registerCallbackDelegate(callbacksDelegate interface{}) (expCallbacks, func()) {
	id, unregister := callbacks.register(callbacksDelegate)
	return expCallbacks{
		usrData:       C.user_callback_data(id),
		statsCallback: C.stats_callback(C.goStatsCallback),
		logCallback:   C.log_callback(C.goLogCallback),
		eventCallback: C.event_callback(C.goEventCallback),
	}, unregister

}
