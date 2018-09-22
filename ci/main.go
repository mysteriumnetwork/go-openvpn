package main

import (
	"fmt"
	"github.com/mysteriumnetwork/go-openvpn/ci/commands"
	"os"
	"strings"
)

func main() {
	help := `
Usage: 
	go run ci/main.go <command>

The commands are:
	test      Runs all unit tests in project
	help      Shows a list of commands
`

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Print(help)
		commands.ExitWithError("Need subcommand as first argument")
	}

	command := strings.ToLower(args[0])
	switch command {
	case "test":
		commands.CommandTestUnit(args[1:])
	case "help":
		fmt.Print(help)
	default:
		fmt.Print(help)
		commands.ExitWithError(fmt.Sprint("Unknown command:", command))
	}
}
