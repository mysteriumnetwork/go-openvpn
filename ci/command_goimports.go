// +build mage

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
