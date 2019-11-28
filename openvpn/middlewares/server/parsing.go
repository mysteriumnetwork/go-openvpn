/*
 * Copyright (C) 2018 The "MysteriumNetwork/go-openvpn" Authors.
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

package server

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	ruleID          = regexp.MustCompile(`^(\d+)$`)
	ruleIDAndKey    = regexp.MustCompile(`^(\d+),(\d+)$`)
	ruleClientEvent = regexp.MustCompile(`^(\w+),(.*)$`)
)

// ParseClientEvent parses OpenVPN management client event.
func ParseClientEvent(line string) (ClientEventType, string, error) {
	match := ruleClientEvent.FindStringSubmatch(line)
	if len(match) < 3 {
		return "", "", errors.New("unable to parse event: " + line)
	}

	event := ClientEventType(match[1])
	return event, match[2], nil
}

// ParseEnvVar parses OpenVPN management client environment variable.
func ParseEnvVar(data string) (string, string, error) {
	slice := strings.SplitN(data, "=", 2)
	if len(slice) == 2 {
		return slice[0], slice[1], nil
	} else if len(slice) == 1 {
		return slice[0], "", nil
	}

	return "", "", errors.New("invalid env var: " + data)
}

// ParseIDAndKey parses CID and KID from OpenVPN management client.
func ParseIDAndKey(data string) (int, int, error) {
	match := ruleIDAndKey.FindStringSubmatch(data)
	if len(match) < 3 {
		return Undefined, Undefined, errors.New("unable to parse identifiers: " + data)
	}

	ID, err := strconv.Atoi(match[1])
	if err != nil {
		return Undefined, Undefined, err
	}

	key, err := strconv.Atoi(match[2])
	if err != nil {
		return Undefined, Undefined, err
	}

	return ID, key, nil
}

// ParseID parses CID from OpenVPN management client.
func ParseID(data string) (int, error) {
	match := ruleID.FindStringSubmatch(data)
	if len(match) < 2 {
		return Undefined, errors.New("unable to parse identifier: " + data)
	}

	ID, err := strconv.Atoi(match[1])
	if err != nil {
		return Undefined, err
	}

	return ID, nil
}
