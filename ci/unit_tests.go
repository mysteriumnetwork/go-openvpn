// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Runs the test suite against the repo
func Test() error {
	mg.Deps(Deps)
	return sh.RunV("go", "test", "-race", "-cover", "../...")
}
