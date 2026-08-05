package main

import (
	"os"

	"devstack/internal/bootstrap"
)

var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {
	runner := bootstrap.NewRunner(version, buildDate)
	exitCode := runner.Run(os.Args)
	os.Exit(exitCode)
}
