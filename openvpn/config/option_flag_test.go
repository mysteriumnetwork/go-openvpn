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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlag_Factory(t *testing.T) {
	option := OptionFlag("enable-something")
	assert.NotNil(t, option)
}

func TestFlag_GetName(t *testing.T) {
	option := OptionFlag("enable-something")
	assert.Equal(t, "enable-something", option.getName())
}

func TestFlag_ToCli(t *testing.T) {
	option := OptionFlag("enable-something")

	optionValue, err := option.toCli()
	assert.NoError(t, err)
	assert.Equal(t, []string{"--enable-something"}, optionValue)
}

func TestFlag_ToFile(t *testing.T) {
	option := OptionFlag("enable-something")

	optionValue, err := option.toFile()
	assert.NoError(t, err)
	assert.Equal(t, "enable-something", optionValue)
}
