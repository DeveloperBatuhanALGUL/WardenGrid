package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/detector"
)

type Reporter struct {
	writer    io.Writer
	jsonLines bool
}

type alertRecord struct {
	Timestamp   string `json:"timestamp"`
	Severity    string `json:"severity"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
	SourceIP    string `json:"source_ip"`
	DestIP      string `json:"dest_ip"`
}

func New(jsonLines bool) *Reporter {
	return &Reporter{
		writer:    os.Stdout,
		jsonLines: jsonLines,
	}
}

func NewWithWriter(w io.Writer, jsonLines bool) *Reporter {
	return &Reporter{
		writer:    w,
		jsonLines: jsonLines,
	}
}

func (r *Reporter) Emit(alert *detector.Alert) error {
	if r.jsonLines {
		return r.emitJSON(alert)
	}
	return r.emitText(alert)
}

func (r *Reporter) emitJSON(alert *detector.Alert) error {
	record := alertRecord{
		Timestamp:   formatTimestamp(alert.Timestamp),
		Severity:    severityLabel(alert.Severity),
		Rule:        alert.Rule,
		Description: alert.Description,
		SourceIP:    alert.SourceIP,
		DestIP:      alert.DestIP,
	}

	encoded, err := json.Marshal(record)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(r.writer, string(encoded))
	return err
}

func (r *Reporter) emitText(alert *detector.Alert) error {
	_, err := fmt.Fprintf(
		r.writer,
		"[%s] %-8s %-28s %s -> %s : %s\n",
		formatTimestamp(alert.Timestamp),
		severityLabel(alert.Severity),
		alert.Rule,
		alert.SourceIP,
		alert.DestIP,
		alert.Description,
	)
	return err
}

func formatTimestamp(ts int64) string {
	if ts == 0 {
		return time.Now().UTC().Format(time.RFC3339)
	}
	return time.Unix(ts, 0).UTC().Format(time.RFC3339)
}

func severityLabel(s detector.Severity) string {
	switch s {
	case detector.SeverityInfo:
		return "INFO"
	case detector.SeverityLow:
		return "LOW"
	case detector.SeverityMedium:
		return "MEDIUM"
	case detector.SeverityHigh:
		return "HIGH"
	case detector.SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}
