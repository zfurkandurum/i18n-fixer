package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

// JSONReporter writes machine-readable JSON output.
type JSONReporter struct{}

func (r *JSONReporter) Report(result *types.AuditResult, w io.Writer) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("writing JSON: %w", err)
	}
	_, err = fmt.Fprintln(w)
	return err
}
