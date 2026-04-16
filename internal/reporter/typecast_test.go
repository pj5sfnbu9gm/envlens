package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/typecast"
)

func sampleTypecastResults() []typecast.Result {
	return []typecast.Result{
		{Key: "PORT", Raw: "8080", Kind: "int", Value: 8080},
		{Key: "DEBUG", Raw: "true", Kind: "bool", Value: true},
		{Key: "RATIO", Raw: "bad", Kind: "float", Error: `cannot cast "bad" to float`},
	}
}

func TestReportTypecast_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultTypecastOptions(&buf)
	if err := ReportTypecast(sampleTypecastResults(), opts); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Error("expected PORT in output")
	}
	if !strings.Contains(out, "int") {
		t.Error("expected int type in output")
	}
	if !strings.Contains(out, "cannot cast") {
		t.Error("expected error message in output")
	}
}

func TestReportTypecast_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultTypecastOptions(&buf)
	opts.Format = "json"
	if err := ReportTypecast(sampleTypecastResults(), opts); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"key"`) {
		t.Error("expected json key field")
	}
	if !strings.Contains(out, `"kind"`) {
		t.Error("expected json kind field")
	}
}

func TestReportTypecast_Empty(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultTypecastOptions(&buf)
	if err := ReportTypecast([]typecast.Result{}, opts); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "KEY") {
		t.Error("expected header even for empty results")
	}
}
