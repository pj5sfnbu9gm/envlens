package reporter

import (
	"encoding/json"
	"io"

	"github.com/user/envlens/internal/differ"
)

// jsonReport is the top-level JSON output structure.
type jsonReport struct {
	From    string        `json:"from"`
	To      string        `json:"to"`
	Changes []jsonChange  `json:"changes"`
	Summary jsonSummary   `json:"summary"`
}

type jsonChange struct {
	Key       string `json:"key"`
	Status    string `json:"status"`
	FromValue string `json:"from_value,omitempty"`
	ToValue   string `json:"to_value,omitempty"`
}

type jsonSummary struct {
	Added     int `json:"added"`
	Removed   int `json:"removed"`
	Changed   int `json:"changed"`
	Unchanged int `json:"unchanged"`
}

func writeJSON(w io.Writer, results []differ.Result, from, to string) error {
	report := jsonReport{
		From:    from,
		To:      to,
		Changes: make([]jsonChange, 0, len(results)),
	}

	for _, r := range results {
		change := jsonChange{
			Key:       r.Key,
			Status:    string(r.Status),
			FromValue: r.FromValue,
			ToValue:   r.ToValue,
		}
		report.Changes = append(report.Changes, change)

		switch r.Status {
		case differ.StatusAdded:
			report.Summary.Added++
		case differ.StatusRemoved:
			report.Summary.Removed++
		case differ.StatusChanged:
			report.Summary.Changed++
		case differ.StatusUnchanged:
			report.Summary.Unchanged++
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
