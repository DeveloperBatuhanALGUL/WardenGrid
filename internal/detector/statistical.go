package detector

import (
	"fmt"
	"math"

	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/protocol/modbus"
)

type intervalSample struct {
	lastTimestamp int64
	intervals     []float64
}

type BaselineDeviationRule struct {
	samples         map[string]*intervalSample
	minSampleSize   int
	zScoreThreshold float64
	minStdDev       float64
}

func NewBaselineDeviationRule() *BaselineDeviationRule {
	return &BaselineDeviationRule{
		samples:         make(map[string]*intervalSample),
		minSampleSize:   5,
		zScoreThreshold: 3.0,
		minStdDev:       0.5,
	}
}

func (r *BaselineDeviationRule) Evaluate(frame *modbus.Frame, d *Detector) *Alert {
	if frame.IsException {
		return nil
	}

	ts := nowOrTimestamp(frame.Timestamp)
	sample, exists := r.samples[frame.SourceIP]
	if !exists {
		r.samples[frame.SourceIP] = &intervalSample{lastTimestamp: ts}
		return nil
	}

	interval := float64(ts - sample.lastTimestamp)
	sample.lastTimestamp = ts
	if interval < 0 {
		interval = 0
	}

	if len(sample.intervals) < r.minSampleSize {
		sample.intervals = append(sample.intervals, interval)
		return nil
	}

	mean, stdDev := meanAndStdDev(sample.intervals)
	sample.intervals = append(sample.intervals[1:], interval)

	if stdDev < r.minStdDev {
		stdDev = r.minStdDev
	}

	zScore := math.Abs(interval-mean) / stdDev
	if zScore < r.zScoreThreshold {
		return nil
	}

	return &Alert{
		Severity:    SeverityLow,
		Rule:        "baseline-timing-deviation",
		Description: fmt.Sprintf("request interval %.2fs deviates %.2f standard deviations from baseline mean %.2fs for %s", interval, zScore, mean, frame.SourceIP),
		SourceIP:    frame.SourceIP,
		DestIP:      frame.DestIP,
		Timestamp:   frame.Timestamp,
	}
}

func meanAndStdDev(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	varianceSum := 0.0
	for _, v := range values {
		diff := v - mean
		varianceSum += diff * diff
	}
	variance := varianceSum / float64(len(values))

	return mean, math.Sqrt(variance)
}
