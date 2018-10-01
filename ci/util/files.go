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
