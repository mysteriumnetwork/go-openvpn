/*
 * Copyright (C) 2019 The "MysteriumNetwork/go-openvpn" Authors.
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
