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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/go-openvpn/ci/util"
)

// Checks if golint exists, if not installs it
func GetLint() error {
	path, _ := util.GetGoBinaryPath("golint")
	if path != "" {
		fmt.Println("Tool 'golint' already installed")
		return nil
	}
	err := sh.RunV("go", "get", "-u", "golang.org/x/lint/golint")
	if err != nil {
		fmt.Println("Could not go get golint")
		return err
	}
	return nil
}

var packageRegexp = regexp.MustCompile(`\.\./(.*)\/.*\.go`)

func getPackageFromGoLintOutput(line string) string {
	results := packageRegexp.FindAllStringSubmatch(line, -1)
	for i := range results {
		return results[i][1]
	}
	return ""
}

func formatAndPrintGoLintOutput(rawGolint string) {
	packageErrorMap := make(map[string][]string, 0)
	separateLines := strings.Split(rawGolint, "\n")

	for i := range separateLines {
		pkg := getPackageFromGoLintOutput(separateLines[i])
		if val, ok := packageErrorMap[pkg]; ok {
			packageErrorMap[pkg] = append(val, separateLines[i])
		} else {
			lines := []string{separateLines[i]}
			packageErrorMap[pkg] = lines
		}
	}

	fmt.Println()
	for k := range packageErrorMap {
		fmt.Println("PACKAGE: ", k)
		fmt.Println()
		for _, v := range packageErrorMap[k] {
			fmt.Println(v)
		}
		fmt.Println()
	}
}

// Reports linting errors in the solution
func GoLint() error {
	mg.Deps(GetLint)
	path, err := util.GetGoBinaryPath("golint")
	if err != nil {
		return err
	}
	var files []string
	err = filepath.Walk("../", func(path string, info os.FileInfo, err error) error {
		if util.IsPathExcluded(path) {
			return nil
		}
		if info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	args := []string{"--set_exit_status", "--min_confidence=1"}
	args = append(args, files...)
	output, err := sh.Output(path, args...)
	exitStatus := sh.ExitStatus(err)
	if exitStatus == 0 {
		fmt.Println("No linting errors")
		return nil
	}

	formatAndPrintGoLintOutput(output)
	fmt.Println("Linting failed!")
	return err
}
