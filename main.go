package main

import (
	"github.com/mysteriumnetwork/openvpnv3-go-bindings/process"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	process.CheckLibrary()

	p := process.NewProcess()

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
