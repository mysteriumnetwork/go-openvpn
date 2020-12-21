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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigToArguments(t *testing.T) {
	config := GenericConfig{}
	config.AddOptions(
		OptionFlag("flag"),
		OptionFlag("spacy flag"),
		OptionParam("value", "1234"),
		OptionParam("very-value", "1234", "5678"),
		OptionParam("spacy value", "1234", "5678"),
	)

	arguments, err := config.ToArguments()
	assert.NoError(t, err)
	assert.Equal(t,
		[]string{
			"--flag",
			"--spacy flag",
			"--value", "1234",
			"--very-value", "1234", "5678",
			"--spacy value", "1234", "5678",
		},
		arguments,
	)
}

func TestSpacedValuesArePassedAsSingleArg(t *testing.T) {
	config := GenericConfig{}
	config.AddOptions(
		OptionParam("value1", "with spaces"),
		OptionFile("value2", "file content", filepath.Join("testdataoutput", "name with spaces.txt")),
	)
	args, err := config.ToArguments()
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]string{
			"--value1", "with spaces",
			"--value2", "testdataoutput/name with spaces.txt",
		},
		args,
	)
}
