package detector

import (
	"testing"

	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/protocol/modbus"
	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/simulator"
)

func TestDetectorNormalTrafficNoAlerts(t *testing.T) {
	d := New()
	cfg := simulator.Config{
		SourceIP:      "10.0.0.5",
		DestIP:        "10.0.0.10",
		UnitID:        1,
		Scenario:      simulator.ScenarioNormalPolling,
		FrameCount:    5,
		BaseTimestamp: 1700000000,
	}

	frames := simulator.Generate(cfg)
	for _, gf := range frames {
		frame, err := modbus.ParseFrame(gf.Raw, gf.SourceIP, gf.DestIP, gf.Timestamp)
		if err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		alerts := d.Process(frame)
		if len(alerts) != 0 {
			t.Errorf("expected no alerts for normal traffic, got %d", len(alerts))
		}
	}
}

func TestDetectorUnknownFunctionCode(t *testing.T) {
	d := New()
	cfg := simulator.Config{
		SourceIP:      "10.0.0.5",
		DestIP:        "10.0.0.10",
		UnitID:        1,
		Scenario:      simulator.ScenarioUnexpectedFunctionCode,
		FrameCount:    5,
		BaseTimestamp: 1700000000,
	}

	frames := simulator.Generate(cfg)
	totalAlerts := 0
	for _, gf := range frames {
		frame, err := modbus.ParseFrame(gf.Raw, gf.SourceIP, gf.DestIP, gf.Timestamp)
		if err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		alerts := d.Process(frame)
		totalAlerts += len(alerts)
	}
	if totalAlerts == 0 {
		t.Errorf("expected at least one alert for unknown function code scenario")
	}
}

func TestDetectorProtectedRegisterWrite(t *testing.T) {
	d := New()
	cfg := simulator.Config{
		SourceIP:      "10.0.0.5",
		DestIP:        "10.0.0.10",
		UnitID:        1,
		Scenario:      simulator.ScenarioUnauthorizedRegisterWrite,
		FrameCount:    4,
		BaseTimestamp: 1700000000,
	}

	frames := simulator.Generate(cfg)
	found := false
	for _, gf := range frames {
		frame, err := modbus.ParseFrame(gf.Raw, gf.SourceIP, gf.DestIP, gf.Timestamp)
		if err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		alerts := d.Process(frame)
		for _, a := range alerts {
			if a.Rule == "protected-register-write" {
				found = true
				if a.Severity != SeverityCritical {
					t.Errorf("expected critical severity, got %v", a.Severity)
				}
			}
		}
	}
	if !found {
		t.Errorf("expected protected-register-write alert")
	}
}

func TestDetectorAbnormalWriteFrequency(t *testing.T) {
	d := New()
	cfg := simulator.Config{
		SourceIP:      "10.0.0.5",
		DestIP:        "10.0.0.10",
		UnitID:        1,
		Scenario:      simulator.ScenarioAbnormalWriteFrequency,
		FrameCount:    10,
		BaseTimestamp: 1700000000,
	}

	frames := simulator.Generate(cfg)
	found := false
	for _, gf := range frames {
		frame, err := modbus.ParseFrame(gf.Raw, gf.SourceIP, gf.DestIP, gf.Timestamp)
		if err != nil {
			t.Fatalf("unexpected parse error: %v", err)
		}
		alerts := d.Process(frame)
		for _, a := range alerts {
			if a.Rule == "abnormal-write-frequency" {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("expected abnormal-write-frequency alert")
	}
}
