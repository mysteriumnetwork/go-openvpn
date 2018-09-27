package util

import (
	"os"
	"path/filepath"
	"strings"
)

var excludedDirs = []string{".git", "vendor"}

// IsPathExcluded determines if the provided path is excluded from common searches
func IsPathExcluded(path string) bool {
	for _, exclude := range excludedDirs {
		if strings.Contains(path, "/"+exclude) {
			return true
		}
	}
	return false
}

// GetProjectFileDirectories returns all the project directories excluding git and vendor
func GetProjectFileDirectories() ([]string, error) {
	directories := make([]string, 0)

	err := filepath.Walk("../", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !IsPathExcluded(path) {
			directories = append(directories, path)
		}
		return nil
	})
	return directories, err
}
