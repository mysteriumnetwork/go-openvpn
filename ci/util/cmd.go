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
	"go/build"
	"os"
	"path"

	"github.com/magefile/mage/sh"
)

// GetGoPath returns the go path
func GetGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

// GetGoBinaryPath looks for the given binary in path, if not checks if it's in $GOPATH/bin
func GetGoBinaryPath(binaryName string) (string, error) {
	res, err := sh.Output("which", binaryName)
	if err == nil {
		return res, nil
	}
	gopath := GetGoPath()
	binaryUnderGopath := path.Join(gopath, "bin", binaryName)
	if _, err := os.Stat(binaryUnderGopath); os.IsNotExist(err) {
		return "", err
	}
	return binaryUnderGopath, nil
}
