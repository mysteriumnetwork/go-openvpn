package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mysteriumnetwork/openvpnv3-go-bindings/process"
)

type callbacks interface {
	process.Logger
	process.EventConsumer
	process.StatsConsumer
}

type loggingCallbacks struct {
}

func (lc *loggingCallbacks) Log(text string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		fmt.Println("Openvpn log >>", line)
	}
}

func (lc *loggingCallbacks) OnEvent(event process.Event) {
	fmt.Printf("Openvpn event >> %+v\n", event)
}

func (lc *loggingCallbacks) OnStats(stats process.Statistics) {
	fmt.Printf("Openvpn stats >> %+v\n", stats)
}

var _ callbacks = &loggingCallbacks{}

type StdoutLogger func(text string)

func (lc StdoutLogger) Log(text string) {
	lc(text)
}

func main() {

	var logger = func(text string) {
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			fmt.Println("Library check >>", line)
		}
	}

	process.CheckLibrary(StdoutLogger(logger))

	p := process.NewProcess(&loggingCallbacks{})

	bytes, err := ioutil.ReadFile("client.ovpn")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	p.RunWithArgs(string(bytes))
	err = p.WaitFor()
	if err != nil {
		fmt.Println("Process error: ", err)
	} else {
		fmt.Println("Graceful exit")
	}

}
