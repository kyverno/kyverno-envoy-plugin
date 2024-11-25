package main

import (
	"os"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/root"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	log.SetLogger(zap.New(zap.UseDevMode(true)))
	root := root.Command()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
