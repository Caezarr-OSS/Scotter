// Scotter is a scaffolding tool for Go projects that allows
// rapid generation of project structures with integrated CI/CD workflows.
package main

import (
	"fmt"
	"os"

	"github.com/caezarr-oss/scotter/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
