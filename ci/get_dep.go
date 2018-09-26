// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Installs the package mangement tool - dep
func Dep() error {
	err := sh.Run("which", "dep")
	if err == nil {
		fmt.Println("Tool 'dep' already installed")
		return nil
	}
	err = sh.RunV("go", "get", "github.com/golang/dep/cmd/dep")
	if err != nil {
		fmt.Println("Could not go get dep")
		return err
	}
	return nil
}
