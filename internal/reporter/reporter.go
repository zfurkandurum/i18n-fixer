package reporter

import (
	"fmt"
	"io"

	"github.com/zfurkandurum/i18n-fixer/internal/types"
)

// Reporter generates output from audit results.
type Reporter interface {
	Report(result *types.AuditResult, w io.Writer) error
}

// New creates a reporter by format name.
func New(format string) (Reporter, error) {
	switch format {
	case "console":
		return &ConsoleReporter{}, nil
	case "json":
		return &JSONReporter{}, nil
	case "prompt":
		return &PromptReporter{}, nil
	default:
		return nil, fmt.Errorf("unknown format: %s (use: console, json, prompt)", format)
	}
}
