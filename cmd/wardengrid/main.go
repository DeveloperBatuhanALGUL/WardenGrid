package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/detector"
	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/protocol/modbus"
	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/report"
	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/simulator"
)

func main() {
	scenarioFlag := flag.String("scenario", "normal", "simulation scenario: normal, unknown-function, write-frequency, protected-write")
	countFlag := flag.Int("count", 10, "number of frames to generate")
	jsonFlag := flag.Bool("json", false, "emit alerts as JSON lines instead of text")
	flag.Parse()

	scenario, err := resolveScenario(*scenarioFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg := simulator.Config{
		SourceIP:      "192.168.10.50",
		DestIP:        "192.168.10.5",
		UnitID:        1,
		Scenario:      scenario,
		FrameCount:    *countFlag,
		BaseTimestamp: 1700000000,
	}

	frames := simulator.Generate(cfg)
	d := detector.New()
	r := report.New(*jsonFlag)

	fmt.Fprintf(os.Stderr, "WardenGrid: processing %d simulated frames under scenario %q\n", len(frames), *scenarioFlag)

	alertCount := 0
	for _, gf := range frames {
		frame, err := modbus.ParseFrame(gf.Raw, gf.SourceIP, gf.DestIP, gf.Timestamp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WardenGrid: failed to parse frame: %v\n", err)
			continue
		}

		alerts := d.Process(frame)
		for _, a := range alerts {
			if err := r.Emit(a); err != nil {
				fmt.Fprintf(os.Stderr, "WardenGrid: failed to emit alert: %v\n", err)
				continue
			}
			alertCount++
		}
	}

	fmt.Fprintf(os.Stderr, "WardenGrid: run complete. %d alert(s) raised.\n", alertCount)
}

func resolveScenario(name string) (simulator.ScenarioType, error) {
	switch name {
	case "normal":
		return simulator.ScenarioNormalPolling, nil
	case "unknown-function":
		return simulator.ScenarioUnexpectedFunctionCode, nil
	case "write-frequency":
		return simulator.ScenarioAbnormalWriteFrequency, nil
	case "protected-write":
		return simulator.ScenarioUnauthorizedRegisterWrite, nil
	default:
		return 0, fmt.Errorf("unknown scenario %q", name)
	}
}
