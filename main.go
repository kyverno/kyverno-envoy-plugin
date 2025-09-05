package main

import (
	"os"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/root"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/logging"
)

func main() {
	setupLogging()
	logger := logging.WithName("kyverno-envoy-plugin")
	root := root.Command()
	if err := root.Execute(); err != nil {
		logger.Error(err, "failed to execute root command")
		os.Exit(1)
	}
}

func setupLogging() {
	if err := logging.Setup(logging.JSONFormat, logging.ISO8601, logging.LogLevel, false); err != nil {
		logging.Error(err, "failed to setup logging")
		os.Exit(1)
	}
}
