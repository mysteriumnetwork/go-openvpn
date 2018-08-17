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
	"strings"
)

//export GoStatsCallback
func GoStatsCallback(ptr C.user_data, stats C.conn_stats) {

	fmt.Printf("%+v\n", stats)
}

//export GoLogCallback
func GoLogCallback(ptr C.user_data, cStr *C.char) {
	goStr := C.GoString(cStr)
	logLines := strings.Split(goStr, "\n")
	for _ , logLine := range logLines {
		fmt.Println("Openvpn >>" , logLine)
	}

}

//export GoEventCallback
func GoEventCallback(ptr C.user_data, event C.conn_event) {
	name := C.GoString(event.name)
	info := C.GoString(event.info)
	fmt.Println("Event >>", event.error, event.fatal, name, info)
}

type Process struct {
	resChan chan error
}

func CheckLibrary() {
	C.checkLibrary(nil, C.log_callback(C.GoLogCallback))
}

func NewProcess() (*Process) {
	return &Process{
		resChan: make(chan error),
	}
}

func (p * Process) RunWithArgs(args... string) {
	go func() {

		cPtr:=NewCharPointer(args[0])
		defer cPtr.Delete()

		res , err := C.initProcess( cPtr.Ptr, C.user_data(nil) , C.stats_callback(C.GoStatsCallback), C.log_callback(C.GoLogCallback), C.event_callback(C.GoEventCallback))
		if err != nil {
			p.resChan <- err
		} else if res != 0 {
			p.resChan <- fmt.Errorf("res error: %v", res)
		} else {
			close(p.resChan)
		}

	}()
}

func (p * Process) WaitFor() error {
	return <- p.resChan
}

func (p * Process) Stop() {

}
