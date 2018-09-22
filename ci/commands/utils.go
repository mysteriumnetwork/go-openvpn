package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
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

// GoPackage creates command to go binary which was installed as Go gettable dependency
func GoPackage(binary string, args ...string) *exec.Cmd {
	return exec.Command(
		path.Join(os.Getenv("GOPATH"), "bin", binary),
		args...,
	)
}

func ExitWithError(message string) {
	fmt.Println(message)
	os.Exit(1)
}
