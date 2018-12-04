package main

import (
	"os"

	"github.com/soushin/bazel-multiprojects/pkg/ops/cmd"
)

func main() {
	if err := cmd.NewDefaultCommand().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
