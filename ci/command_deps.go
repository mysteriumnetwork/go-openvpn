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

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/go-openvpn/ci/util"
)

// Installs go dependencies
func Deps() error {
	mg.Deps(Dep)
	dir, err := util.GetGoBinaryPath("dep")
	if err != nil {
		fmt.Println("Could not find dep")
		return err
	}
	fmt.Println("Installing dependencies with", dir)
	err = sh.RunV(dir, "ensure", "-v")
	if err != nil {
		fmt.Println("Installing dependencies FAIL!")
	} else {
		fmt.Println("Installing dependencies DONE!")
	}
	return err
}
