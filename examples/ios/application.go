package ios

import (
	"fmt"
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

func (lc *loggingCallbacks) OnEvent(event openvpn3.Event) {
	fmt.Printf("Openvpn event >> %+v\n", event)
}

func (lc *loggingCallbacks) OnStats(stats openvpn3.Statistics) {
	fmt.Printf("Openvpn stats >> %+v\n", stats)
}

var _ callbacks = &loggingCallbacks{}

type StdoutLogger func(text string)

func (lc StdoutLogger) Log(text string) {
	lc(text)
}

func StartSession() {

	var logger StdoutLogger = func(text string) {
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			fmt.Println("Library check >>", line)
		}
	}
	openvpn3.SelfCheck(logger)

	profile := ""
	profileCreds := openvpn3.Credentials{
		Username: "abc",
		Password: "def",
	}

	session := openvpn3.NewMobileSession(&loggingCallbacks{}, &openvpn3.NoOpTunnelSetup{})
	session.Start(profile, profileCreds)

	err := session.Wait()
	if err != nil {
		fmt.Println("Openvpn3 error: ", err)
	} else {
		fmt.Println("Graceful exit")
	}

}
