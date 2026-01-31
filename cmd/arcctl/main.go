// Package main provides the arcctl CLI entry point.
package main

import (
	"fmt"
	"os"

	"github.com/architect-io/arcctl/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
