package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/zfurkandurum/i18n-fixer/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.SetVersionInfo(version, commit, date)
	if err := cli.Execute(); err != nil {
		var issuesErr *cli.IssuesFoundError
		if errors.As(err, &issuesErr) {
			// Issues found — not a program error, just exit 1 for CI
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
}
