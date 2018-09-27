// +build mage

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
