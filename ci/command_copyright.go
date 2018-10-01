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
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var copyrightRegex = regexp.MustCompile(`Copyright \(C\) \d{4} The "MysteriumNetwork/go-openvpn"`)

func getFilesWithoutCopyright() ([]string, error) {
	badFiles := make([]string, 0)
	var excludedDirs = []string{".git", "vendor"}
	err := filepath.Walk("../", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, exclude := range excludedDirs {
			if strings.Contains(path, "/"+exclude) {
				return nil
			}
		}
		extension := filepath.Ext(path)
		if extension != ".go" {
			return nil
		}
		contents, err := ioutil.ReadFile(path)
		match := copyrightRegex.Match(contents)
		if !match {
			badFiles = append(badFiles, path)
		}
		return nil
	})
	return badFiles, err
}

// Checks for copyright headers in files
func Copyright() error {
	badFiles, err := getFilesWithoutCopyright()
	if err != nil {
		return err
	}
	if len(badFiles) != 0 {
		fmt.Println("Following files are missing copyright headers:")
		for _, v := range badFiles {
			fmt.Println(v)
		}
		return errors.New("Missing copyright headers")
	}
	fmt.Println("All files have required copyright headers!")
	return nil
}
