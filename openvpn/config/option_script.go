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
	"fmt"
	"path/filepath"
)

// Script interface provides single method - for given script dir, provide full script path
type Script interface {
	FullPath(scriptDir string) string
}

// SimplePath function constructs unquoted script path from given name
func SimplePath(name string) Script {
	return &scriptPath{
		name,
	}
}

// QuotedPath function constructs quoted script path from given name
func QuotedPath(name string) Script {
	return &quotedScriptPath{
		scriptPath{
			name,
		},
	}
}

type scriptPath struct {
	scriptName string
}

func (sp *scriptPath) FullPath(scriptDir string) string {
	return filepath.Join(scriptDir, sp.scriptName)
}

type quotedScriptPath struct {
	scriptPath
}

func (sp *quotedScriptPath) FullPath(scriptDir string) string {
	return wrapWithDoubleQuotes(sp.scriptPath.FullPath(scriptDir))
}

func wrapWithDoubleQuotes(val string) string {
	return fmt.Sprintf(`"%s"`, val)
}
