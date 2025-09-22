package version

import (
	"fmt"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/version"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:          "version",
		Short:        "Print the version informations",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Version: %s\n", version.Version()); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Time: %s\n", version.Time()); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Git commit ID: %s\n", version.Hash()); err != nil {
				return err
			}
			return nil
		},
	}
}
