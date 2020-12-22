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

package config

// OptionFlag creates flag style option
func OptionFlag(name string) optionFlag {
	return optionFlag{name}
}

type optionFlag struct {
	name string
}

func (option optionFlag) getName() string {
	return option.name
}

func (option optionFlag) toCli() ([]string, error) {
	return []string{"--" + option.name}, nil
}

func (option optionFlag) toFile() (string, error) {
	return option.name, nil
}
