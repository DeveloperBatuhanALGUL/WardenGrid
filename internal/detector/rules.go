package detector

import (
	"fmt"

	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/protocol/modbus"
)

type UnknownFunctionCodeRule struct{}

func (r *UnknownFunctionCodeRule) Evaluate(frame *modbus.Frame, d *Detector) *Alert {
	if frame.IsException {
		return nil
	}
	if frame.Function.IsKnown() {
		return nil
	}
	return &Alert{
		Severity:    SeverityHigh,
		Rule:        "unknown-function-code",
		Description: fmt.Sprintf("unrecognized Modbus function code 0x%X from %s", byte(frame.Function), frame.SourceIP),
		SourceIP:    frame.SourceIP,
		DestIP:      frame.DestIP,
		Timestamp:   frame.Timestamp,
	}
}

type ProtectedRegisterWriteRule struct{}

func (r *ProtectedRegisterWriteRule) Evaluate(frame *modbus.Frame, d *Detector) *Alert {
	if frame.IsException {
		return nil
	}
	if !frame.Function.IsWrite() {
		return nil
	}
	if !d.isProtectedAddress(frame.StartAddress) {
		return nil
	}
	return &Alert{
		Severity:    SeverityCritical,
		Rule:        "protected-register-write",
		Description: fmt.Sprintf("write to protected register 0x%X by %s", frame.StartAddress, frame.SourceIP),
		SourceIP:    frame.SourceIP,
		DestIP:      frame.DestIP,
		Timestamp:   frame.Timestamp,
	}
}

type WriteFrequencyRule struct{}

func (r *WriteFrequencyRule) Evaluate(frame *modbus.Frame, d *Detector) *Alert {
	if frame.IsException {
		return nil
	}
	if !frame.Function.IsWrite() {
		return nil
	}

	ts := nowOrTimestamp(frame.Timestamp)
	count := d.recordWrite(frame.SourceIP, ts)

	if count <= d.writeThreshold {
		return nil
	}

	return &Alert{
		Severity:    SeverityMedium,
		Rule:        "abnormal-write-frequency",
		Description: fmt.Sprintf("%d write requests from %s within %d seconds", count, frame.SourceIP, d.writeWindowSecs),
		SourceIP:    frame.SourceIP,
		DestIP:      frame.DestIP,
		Timestamp:   frame.Timestamp,
	}
}
