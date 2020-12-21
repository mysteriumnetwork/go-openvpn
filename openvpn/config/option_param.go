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

import "strings"

// OptionParam creates --name value1 value2 style config parameter
func OptionParam(name string, values ...string) optionParam {
	return optionParam{name, values}
}

type optionParam struct {
	name   string
	values []string
}

func (option optionParam) getName() string {
	return option.name
}

func (option optionParam) toCli() ([]string, error) {
	return append([]string{"--" + option.name}, option.values...), nil
}

func (option optionParam) toFile() (string, error) {
	return option.name + " " + strings.Join(option.values, " "), nil
}
