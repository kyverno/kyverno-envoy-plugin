package main

import (
	"os"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/root"
)

func main() {
	root := root.Command()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
