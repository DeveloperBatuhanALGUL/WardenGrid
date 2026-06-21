package detector

import (
	"time"

	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/protocol/modbus"
)

type Severity int

const (
	SeverityInfo Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

type Alert struct {
	Severity    Severity
	Rule        string
	Description string
	SourceIP    string
	DestIP      string
	Timestamp   int64
}

type Detector struct {
	rules           []Rule
	writeHistory    map[string][]int64
	writeWindowSecs int64
	writeThreshold  int
	protectedRanges []AddressRange
}

type AddressRange struct {
	Start uint16
	End   uint16
}

type Rule interface {
	Evaluate(frame *modbus.Frame, d *Detector) *Alert
}

func New() *Detector {
	return &Detector{
		writeHistory:    make(map[string][]int64),
		writeWindowSecs: 10,
		writeThreshold:  5,
		protectedRanges: []AddressRange{
			{Start: 0xFF00, End: 0xFFFF},
		},
		rules: []Rule{
			&UnknownFunctionCodeRule{},
			&ProtectedRegisterWriteRule{},
			&WriteFrequencyRule{},
			NewBaselineDeviationRule(),
		},
	}
}

func (d *Detector) Process(frame *modbus.Frame) []*Alert {
	alerts := make([]*Alert, 0)
	for _, rule := range d.rules {
		if alert := rule.Evaluate(frame, d); alert != nil {
			alerts = append(alerts, alert)
		}
	}
	return alerts
}

func (d *Detector) recordWrite(sourceIP string, timestamp int64) int {
	history := d.writeHistory[sourceIP]
	cutoff := timestamp - d.writeWindowSecs
	filtered := history[:0]
	for _, ts := range history {
		if ts >= cutoff {
			filtered = append(filtered, ts)
		}
	}
	filtered = append(filtered, timestamp)
	d.writeHistory[sourceIP] = filtered
	return len(filtered)
}

func (d *Detector) isProtectedAddress(addr uint16) bool {
	for _, r := range d.protectedRanges {
		if addr >= r.Start && addr <= r.End {
			return true
		}
	}
	return false
}

func nowOrTimestamp(ts int64) int64 {
	if ts == 0 {
		return time.Now().Unix()
	}
	return ts
}
