// +build mage

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
