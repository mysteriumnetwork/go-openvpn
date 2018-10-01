// +build mage

/*
 * Copyright (C) 2018 The "MysteriumNetwork/go-openvpn" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/go-openvpn/ci/util"
)

// Fetches the goimports binary
func GetImports() error {
	path, _ := util.GetGoBinaryPath("goimports")
	if path != "" {
		fmt.Println("Tool 'goimports' already installed")
		return nil
	}
	err := sh.RunV("go", "get", "golang.org/x/tools/cmd/goimports")
	if err != nil {
		fmt.Println("Could not go get goimports")
		return err
	}
	return nil
}

// Checks for issues with go imports
func GoImports() error {
	mg.Deps(GetImports)
	path, err := util.GetGoBinaryPath("goimports")
	if err != nil {
		fmt.Println("Tool 'goimports' not found")
		return err
	}
	args := []string{"-e", "-l"}
	dirsToLook := make([]string, 0)
	res, _ := ioutil.ReadDir("../")
	for _, v := range res {
		if !v.IsDir() {
			extension := filepath.Ext(v.Name())
			if extension != ".go" {
				continue
			}
		}
		path := "../" + v.Name()
		if !util.IsPathExcluded(path) {
			dirsToLook = append(dirsToLook, path)
		}
	}
	args = append(args, dirsToLook...)
	out, err := sh.Output(path, args...)
	if err != nil {
		fmt.Println("Could not run goimports")
		return err
	}
	if len(out) != 0 {
		fmt.Println("The following files contain go import errors:")
		fmt.Println(out)
		return errors.New("Not all imports follow the goimports format.")
	}
	fmt.Println("Goimports is happy - all files are OK!")
	return nil
}
