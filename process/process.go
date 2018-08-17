package process
/*
#include "process.h"
#include <stdio.h>

extern void GoStatsCallback(conn_stats stats);
*/
import "C"
import (
	"fmt"
)

//export GoStatsCallback
func GoStatsCallback(stats C.conn_stats) {
	fmt.Println("Golang callback!")
	fmt.Printf("bytes: %v\n", stats.bytes_out)
	fmt.Printf("%+v\n", stats)
}


type Process struct {
	resChan chan error
}

func NewProcess() (*Process) {
	return &Process{
		resChan: make(chan error),
	}
}

func (p * Process) RunWithArgs(args... string) {
	go func() {

		procArgs := &Args{}
		procArgs.AddAll(args...)
		defer procArgs.Free()

		res , err := C.initProcess( procArgs.cPointer(), procArgs.cCount(), C.stats_callback(C.GoStatsCallback))
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
