package simulator

import (
	"testing"

	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/protocol/modbus"
)

func TestGenerateNormalPolling(t *testing.T) {
	cfg := Config{
		SourceIP:      "10.0.0.5",
		DestIP:        "10.0.0.10",
		UnitID:        1,
		Scenario:      ScenarioNormalPolling,
		FrameCount:    5,
		BaseTimestamp: 1700000000,
	}

	frames := Generate(cfg)
	if len(frames) != 5 {
		t.Fatalf("expected 5 frames, got %d", len(frames))
	}

	for _, gf := range frames {
		frame, err := modbus.ParseFrame(gf.Raw, gf.SourceIP, gf.DestIP, gf.Timestamp)
		if err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		if frame.Function != modbus.FuncReadHoldingRegisters {
			t.Errorf("expected read holding registers, got %v", frame.Function)
		}
	}
}

func TestGenerateUnexpectedFunctionCode(t *testing.T) {
	cfg := Config{
		SourceIP:      "10.0.0.5",
		DestIP:        "10.0.0.10",
		UnitID:        1,
		Scenario:      ScenarioUnexpectedFunctionCode,
		FrameCount:    5,
		BaseTimestamp: 1700000000,
	}

	frames := Generate(cfg)
	found := false
	for _, gf := range frames {
		frame, err := modbus.ParseFrame(gf.Raw, gf.SourceIP, gf.DestIP, gf.Timestamp)
		if err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		if !frame.Function.IsKnown() {
			found = true
		}
	}
	if !found {
		t.Errorf("expected at least one unknown function code frame")
	}
}

func TestGenerateUnauthorizedRegisterWrite(t *testing.T) {
	cfg := Config{
		SourceIP:      "10.0.0.5",
		DestIP:        "10.0.0.10",
		UnitID:        1,
		Scenario:      ScenarioUnauthorizedRegisterWrite,
		FrameCount:    4,
		BaseTimestamp: 1700000000,
	}

	frames := Generate(cfg)
	last := frames[len(frames)-1]
	frame, err := modbus.ParseFrame(last.Raw, last.SourceIP, last.DestIP, last.Timestamp)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if frame.StartAddress != 0xFF00 {
		t.Errorf("expected start address 0xFF00, got 0x%X", frame.StartAddress)
	}
}
