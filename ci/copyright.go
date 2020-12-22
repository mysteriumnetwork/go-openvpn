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

package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"

	"github.com/mysteriumnetwork/go-ci/util"
)

var copyrightRegex = regexp.MustCompile(`Copyright \(C\) \d{4} BlockDev AG`)

func getFilesWithoutCopyright(dirsToCheck []string) ([]string, error) {
	badFiles := make([]string, 0)
	gopath := util.GetGoPath()

	for i := range dirsToCheck {
		absolutePath := path.Join(gopath, "src", dirsToCheck[i])
		files, err := ioutil.ReadDir(absolutePath)
		if err != nil {
			return badFiles, err
		}
		for j := range files {
			if files[j].IsDir() {
				continue
			}
			extension := filepath.Ext(files[j].Name())
			if extension != ".go" {
				continue
			}
			contents, err := ioutil.ReadFile(path.Join(absolutePath, files[j].Name()))
			if err != nil {
				return nil, err
			}
			match := copyrightRegex.Match(contents)
			if !match {
				badFiles = append(badFiles, path.Join(dirsToCheck[i], files[j].Name()))
			}
		}
	}
	return badFiles, nil
}

// Copyright checks for copyright headers in files
func Copyright(path string, excludes ...string) error {
	res, err := util.GetPackagePathsWithExcludes(path, excludes...)
	if err != nil {
		fmt.Println("go list crashed")
		return err
	}
	badFiles, err := getFilesWithoutCopyright(res)
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

// CopyrightD checks for copyright headers in files
//
// Instead of packages, it operates on directories, thus it is compatible with gomodules outside GOPATH.
//
// Example:
//     commands.CopyrightD(".", "docs")
func CopyrightD(path string, excludes ...string) error {
	var allExcludes []string
	allExcludes = append(allExcludes, excludes...)
	allExcludes = append(allExcludes, util.GoLintExcludes()...)
	res, err := util.GetProjectFileDirectories(allExcludes)
	if err != nil {
		fmt.Println("copyright: go list crashed")
		return err
	}
	badFiles, err := getFilesWithoutCopyrightD(res)
	if err != nil {
		fmt.Println("copyright: error listing files")
		return err
	}
	if len(badFiles) != 0 {
		fmt.Println("copyright: following files are missing copyright headers:")
		for _, v := range badFiles {
			fmt.Println(v)
		}
		return errors.New("copyright: missing copyright headers")
	}
	fmt.Println("copyright: all files have required copyright headers!")
	return nil
}

func getFilesWithoutCopyrightD(dirsToCheck []string) ([]string, error) {
	badFiles := make([]string, 0)

	for i := range dirsToCheck {
		files, err := ioutil.ReadDir(dirsToCheck[i])
		if err != nil {
			return badFiles, err
		}
		for j := range files {
			if files[j].IsDir() {
				continue
			}
			extension := filepath.Ext(files[j].Name())
			if extension != ".go" {
				continue
			}
			contents, err := ioutil.ReadFile(path.Join(dirsToCheck[i], files[j].Name()))
			if err != nil {
				return nil, err
			}
			match := copyrightRegex.Match(contents)
			if !match {
				badFiles = append(badFiles, path.Join(dirsToCheck[i], files[j].Name()))
			}
		}
	}
	return badFiles, nil
}
