package cmd

import (
	"flag"
	"os"

	"github.com/soushin/bazel-multiprojects/pkg/ops/cmd/deploy"
	"github.com/spf13/cobra"
)

func NewDefaultCommand() *cobra.Command {
	stdOut := os.Stdout

	c := &cobra.Command{
		Use:   "ops",
		Short: "Ops is command for operating service",
		Long:  "",
	}

	c.AddCommand(
		deploy.NewCmdDeploy(stdOut),
	)
	c.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	flag.CommandLine.Parse([]string{})
	return c
}
