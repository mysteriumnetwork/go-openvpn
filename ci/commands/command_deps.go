package commands

import (
	"fmt"
	"os/exec"
)

// CommandDependencies installs all project dependencies
func CommandDependencies(_ []string) {
	if err := exec.Command("which", "dep").Run(); err == nil {
		fmt.Println("Tool 'dep' already installed")
		return
	}

	MustRun("go", "get", "github.com/golang/dep/cmd/dep")
	MustRunCommand(GoPackage("dep", "ensure"))
}
