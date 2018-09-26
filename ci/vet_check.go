// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Checks that the source is compliant with go vet
func Vet() error {
	out, err := sh.Output("go", "vet", "-composites=false", "../...")
	fmt.Print(out)
	if err != nil {
		return err
	}
	fmt.Println("All files are compliant with go vet")
	return nil
}
