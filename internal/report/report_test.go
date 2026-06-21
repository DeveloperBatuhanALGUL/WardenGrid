package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/detector"
)

func TestEmitTextFormat(t *testing.T) {
	var buf bytes.Buffer
	r := NewWithWriter(&buf, false)

	alert := &detector.Alert{
		Severity:    detector.SeverityCritical,
		Rule:        "protected-register-write",
		Description: "write to protected register 0xFF00 by 10.0.0.5",
		SourceIP:    "10.0.0.5",
		DestIP:      "10.0.0.10",
		Timestamp:   1700000000,
	}

	if err := r.Emit(alert); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "CRITICAL") {
		t.Errorf("expected output to contain CRITICAL, got: %s", output)
	}
	if !strings.Contains(output, "protected-register-write") {
		t.Errorf("expected output to contain rule name, got: %s", output)
	}
	if !strings.Contains(output, "10.0.0.5") {
		t.Errorf("expected output to contain source IP, got: %s", output)
	}
}

func TestEmitJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := NewWithWriter(&buf, true)

	alert := &detector.Alert{
		Severity:    detector.SeverityHigh,
		Rule:        "unknown-function-code",
		Description: "unrecognized function code",
		SourceIP:    "10.0.0.5",
		DestIP:      "10.0.0.10",
		Timestamp:   1700000000,
	}

	if err := r.Emit(alert); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded alertRecord
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to decode JSON output: %v", err)
	}

	if decoded.Severity != "HIGH" {
		t.Errorf("expected severity HIGH, got %s", decoded.Severity)
	}
	if decoded.Rule != "unknown-function-code" {
		t.Errorf("expected rule unknown-function-code, got %s", decoded.Rule)
	}
}

func TestSeverityLabelAllValues(t *testing.T) {
	cases := map[detector.Severity]string{
		detector.SeverityInfo:     "INFO",
		detector.SeverityLow:      "LOW",
		detector.SeverityMedium:   "MEDIUM",
		detector.SeverityHigh:     "HIGH",
		detector.SeverityCritical: "CRITICAL",
	}

	for sev, expected := range cases {
		if got := severityLabel(sev); got != expected {
			t.Errorf("expected %s for severity %v, got %s", expected, sev, got)
		}
	}
}
