package commands

import (
	"github.com/magefile/mage/sh"
)

// RunTests runs the tests
func RunTests() error {
	err := sh.RunV("go", "test", "../...")
	return err
}
