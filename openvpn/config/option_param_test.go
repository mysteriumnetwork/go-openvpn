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

func TestParam_Factory(t *testing.T) {
	option := OptionParam("very-value", "1234")
	assert.NotNil(t, option)
}

func TestParam_GetName(t *testing.T) {
	option := OptionParam("very-value", "1234")
	assert.Equal(t, "very-value", option.getName())
}

func TestParam_ToCli(t *testing.T) {
	option := OptionParam("very-value", "1234")

	optionValue, err := option.toCli()
	assert.NoError(t, err)
	assert.Equal(t, []string{"--very-value", "1234"}, optionValue)
}

func TestParam_ToFile(t *testing.T) {
	option := OptionParam("very-value", "1234")

	optionValue, err := option.toFile()
	assert.NoError(t, err)
	assert.Equal(t, "very-value 1234", optionValue)
}
