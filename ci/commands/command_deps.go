package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/magefile/mage/sh"
)

// InstallDep installs dep
func InstallDep() error {
	err := sh.Run("which", "dep")
	if err == nil {
		fmt.Println("Tool 'dep' already installed")
		return nil
	}
	err = sh.RunV("go", "get", "github.com/golang/dep/cmd/dep")
	if err != nil {
		return err
	}
	return nil
}

// InstallDependencies installs go dependencies
func InstallDependencies() error {
	depDir := path.Join(os.Getenv("GOPATH"), "bin", "dep")
	err := sh.RunV(depDir, "ensure")
	if err != nil {
		fmt.Println("Failed to find dep under GOPATH, trying with go directory...")
		res, err := sh.Output("which", "go")
		if err != nil {
			return err
		}
		// replace the go part of path with dep
		depDir = res[:len(res)-2] + "dep"
	}

	fmt.Println("Installing dependencies with", depDir)
	err = sh.RunV(depDir, "ensure", "-v")
	if err != nil {
		fmt.Println("Installing dependencies FAIL!")
	} else {
		fmt.Println("Installing dependencies DONE!")
	}
	return err
}
