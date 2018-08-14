package main

import (
	"github.com/mysteriumnetwork/openvpnv3-go-bindings/process"
	"fmt"
)

func main() {
	p := process.NewProcess()
	p.RunWithArgs("abc","labas")
	err := p.WaitFor()
	if err != nil {
		fmt.Println("Process error: ", err)
	}

}
