// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/mysteriumnetwork/go-openvpn/ci/commands"
)

// Installs the required go dependencies
func Deps() error {
	mg.Deps(Dep)
	return commands.InstallDependencies()
}

// Installs the package mangement tool - dep
func Dep() error {
	return commands.InstallDep()
}

// Runs the test suite against the repo
func Test() error {
	return commands.RunTests()
}
