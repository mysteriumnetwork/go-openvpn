package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mysteriumnetwork/go-openvpn/openvpn3"
)

type callbacks interface {
	openvpn3.Logger
	openvpn3.EventConsumer
	openvpn3.StatsConsumer
}

type loggingCallbacks struct {
}

func (lc *loggingCallbacks) Log(text string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		fmt.Println("Openvpn log >>", line)
	}
}

func (lc *loggingCallbacks) OnEvent(event *openvpn3.Event) {
	fmt.Printf("Openvpn event >> %+v\n", event)
}

func (lc *loggingCallbacks) OnStats(stats *openvpn3.Statistics) {
	fmt.Printf("Openvpn stats >> %+v\n", stats)
}

var _ callbacks = &loggingCallbacks{}

type StdoutLogger func(text string)

func (lc StdoutLogger) Log(text string) {
	lc(text)
}

func main() {

	profileName := os.Args[1]

	var logger StdoutLogger = func(text string) {
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			fmt.Println("Library check >>", line)
		}
	}

	openvpn3.SelfCheck(logger)

	session := openvpn3.NewSession(&loggingCallbacks{})

	bytes, err := ioutil.ReadFile(profileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	creds := openvpn3.Credentials{
		Username: "abc",
		Password: "def",
	}

	session.Start(string(bytes), &creds)
	err = session.Wait()
	if err != nil {
		fmt.Println("Openvpn3 error: ", err)
	} else {
		fmt.Println("Graceful exit")
	}

}
