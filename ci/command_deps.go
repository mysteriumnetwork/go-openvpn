// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/go-openvpn/ci/util"
)

// Installs go dependencies
func Deps() error {
	mg.Deps(Dep)
	dir, err := util.GetGoBinaryPath("dep")
	if err != nil {
		fmt.Println("Could not find dep")
		return err
	}
	fmt.Println("Installing dependencies with", dir)
	err = sh.RunV(dir, "ensure", "-v")
	if err != nil {
		fmt.Println("Installing dependencies FAIL!")
	} else {
		fmt.Println("Installing dependencies DONE!")
	}
	return err
}
