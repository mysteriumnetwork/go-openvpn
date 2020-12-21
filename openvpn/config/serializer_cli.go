/*
 * go-openvpn -- Go gettable library for wrapping Openvpn functionality in go way.
 *
 * Copyright (C) 2020 The "MysteriumNetwork/go-openvpn" Authors..
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

package config

import (
	"fmt"
)

// ToArguments serializes openvpn configuration structure to a list of command line arguments which can be passed to openvpn process
// it also serialize file style options to appropriate files inside given config directory
func (config GenericConfig) ToArguments() ([]string, error) {
	arguments := make([]string, 0)

	for _, item := range config.options {
		option, ok := item.(optionCliSerializable)
		if !ok {
			return nil, fmt.Errorf("unserializable option '%s': %#v", item.getName(), item)
		}

		optionValues, err := option.toCli()
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, optionValues...)
	}

	return arguments, nil
}

type optionCliSerializable interface {
	toCli() ([]string, error)
}
