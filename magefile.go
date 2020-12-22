// +build mage

/*
 * go-openvpn -- Go gettable library for wrapping Openvpn functionality in go way.
 *
 * Copyright (C) 2020 BlockDev AG.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License Version 3
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program in the COPYING file.
 * If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/go-ci/commands"
	cicommands "github.com/mysteriumnetwork/go-openvpn/ci"
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
	return cicommands.Copyright("./...", "docs")
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
