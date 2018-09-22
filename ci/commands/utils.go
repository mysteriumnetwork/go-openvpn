package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// MustRunCommand executes the given command and exits the host process for any error
func MustRunCommand(cmd *exec.Cmd) {
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// MustRun executes the given string as command (see MustRunCommand)
func MustRun(cmd string, args ...string) {
	command := exec.Command(cmd, args...)
	MustRunCommand(command)
}

func ExitWithError(message string) {
	fmt.Println(message)
	os.Exit(1)
}
