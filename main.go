package main

import (
	"fmt"
	"os"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/root"
)

func main() {
	root := root.Command()
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
