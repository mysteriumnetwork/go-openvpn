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

package log

import (
	"log"
)

type stdLogger struct {
}

func (l stdLogger) Error(args ...interface{}) {
	var msg []interface{}
	msg = append(msg, "ERROR")
	msg = append(msg, args...)
	log.Println(msg...)
}

func (l stdLogger) Warn(args ...interface{}) {
	var msg []interface{}
	msg = append(msg, "WARN ")
	msg = append(msg, args...)
	log.Println(msg...)
}

func (l stdLogger) Info(args ...interface{}) {
	var msg []interface{}
	msg = append(msg, "INFO ")
	msg = append(msg, args...)
	log.Println(msg...)
}

func (l stdLogger) Debug(args ...interface{}) {
	var msg []interface{}
	msg = append(msg, "DEBUG")
	msg = append(msg, args...)
	log.Println(msg...)
}

func (l stdLogger) Trace(args ...interface{}) {
	var msg []interface{}
	msg = append(msg, "TRACE")
	msg = append(msg, args...)
	log.Println(msg...)
}
