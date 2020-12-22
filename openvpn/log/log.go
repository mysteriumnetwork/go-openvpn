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

package log

// Logger interface to go-openvpn library logger.
// Stdout implementation is used by default.
type Logger interface {
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
}

var logger Logger = stdLogger{}

// UseLogger sets go-openvpn library logger.
func UseLogger(l Logger) {
	logger = l
}

// UseDefaultLogger resets logger to the default logger.
func UseDefaultLogger() {
	logger = stdLogger{}
}

// Error prints an error.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Warn prints a warning.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Info prints information.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Debug prints debug message.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Trace prints trace message.
func Trace(args ...interface{}) {
	logger.Trace(args...)
}
