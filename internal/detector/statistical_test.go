package detector

import "testing"

func TestBaselineDeviationRuleLearnsAndDetects(t *testing.T) {
	rule := NewBaselineDeviationRule()
	sourceIP := "10.0.0.5"

	baseTimestamps := []int64{1000, 1002, 1004, 1006, 1008, 1010, 1012}
	for i, ts := range baseTimestamps {
		frame := mockFrame(sourceIP, ts)
		alert := rule.Evaluate(&frame, nil)
		if i < len(baseTimestamps)-1 && alert != nil {
			t.Errorf("did not expect alert during baseline learning at step %d, got: %v", i, alert.Description)
		}
	}

	spike := mockFrame(sourceIP, 1012+50)
	alert := rule.Evaluate(&spike, nil)
	if alert == nil {
		t.Fatalf("expected alert for timing spike, got nil")
	}
	if alert.Rule != "baseline-timing-deviation" {
		t.Errorf("expected rule baseline-timing-deviation, got %s", alert.Rule)
	}
}

func TestMeanAndStdDev(t *testing.T) {
	values := []float64{2, 2, 2, 2, 2}
	mean, stdDev := meanAndStdDev(values)
	if mean != 2 {
		t.Errorf("expected mean 2, got %f", mean)
	}
	if stdDev != 0 {
		t.Errorf("expected stddev 0, got %f", stdDev)
	}

	values2 := []float64{1, 2, 3, 4, 5}
	mean2, stdDev2 := meanAndStdDev(values2)
	if mean2 != 3 {
		t.Errorf("expected mean 3, got %f", mean2)
	}
	if stdDev2 <= 0 {
		t.Errorf("expected positive stddev, got %f", stdDev2)
	}
}
