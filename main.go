package main

import (
	"os"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/commands/root"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	root := root.Command()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
