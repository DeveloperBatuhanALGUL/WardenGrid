//go:build !windows

package capture

import "testing"

func TestPcapSourceReadBeforeOpenReturnsError(t *testing.T) {
	source := NewPcapSource("nonexistent0")

	_, err := source.Read()
	if err != ErrSourceNotOpen {
		t.Fatalf("expected ErrSourceNotOpen, got %v", err)
	}
}

func TestPcapSourceCloseWithoutOpenIsSafe(t *testing.T) {
	source := NewPcapSource("nonexistent0")

	if err := source.Close(); err != nil {
		t.Fatalf("expected no error closing unopened source, got %v", err)
	}
}

func TestNewPcapSourceDefaults(t *testing.T) {
	source := NewPcapSource("eth0")

	if source.SnapshotLen != 65535 {
		t.Errorf("expected default snapshot length 65535, got %d", source.SnapshotLen)
	}
	if source.Promiscuous {
		t.Errorf("expected promiscuous mode to default to false")
	}
}
