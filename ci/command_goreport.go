// +build mage

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

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func updateReport() error {
	url := "https://goreportcard.com/checks"
	payload := strings.NewReader("repo=github.com%2Fmysteriumnetwork%2Fnode")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Minute * 2,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Goreports responded with status: %v", res.StatusCode)
	}
	return nil
}

// Updates the go report for the repo
func GoReport() error {
	err := updateReport()
	if err != nil {
		fmt.Println("Report update failure")
		return err
	}
	fmt.Println("Report updated")
	return nil
}
