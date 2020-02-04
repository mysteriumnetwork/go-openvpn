// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/go-ci/commands"
)

const buildPath = "./build/morqa"

// Builds the application
func Build() error {
	return sh.RunV("go", "build", "-o", buildPath, "./cmd/main.go")
}

// Run the application
func Run() error {
	return sh.RunV(buildPath,
		"--bind-addr=:8000",
	)
}

// Runs the test suite against the repo
func Test() error {
	return commands.Test("./...")
}

// Report generates goreport
func GoReport() error {
	return commands.GoReport("github.com/mysteriumnetwork/go-openvpn")
}

// Checks for issues with copyrights
func CheckCopyright() error {
	return commands.Copyright("./...", "docs")
}

// Checks for issues with go imports
func CheckGoImports() error {
	return commands.GoImports("./...", "docs")
}

// Reports linting errors in the solution
func CheckGoLint() error {
	return commands.GoLint("./...", "docs")
}

// Checks that the source is compliant with go vet
func CheckGoVet() error {
	return commands.GoVet("./...")
}

// Checks that the source is compliant with go vet
func Check() {
	mg.Deps(CheckGoImports, CheckGoLint, CheckGoVet, CheckCopyright)
}
