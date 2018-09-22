package commands

import "fmt"

// CommandTestUnit runs all unit tests in project
func CommandTestUnit(_ []string) {
	MustRun("go", "test", "./...")

	fmt.Printf("All tests passed.")
}
