//go:build windows

package capture

import "testing"

func TestPcapSourceOpenReturnsUnsupportedOnWindows(t *testing.T) {
	source := NewPcapSource("nonexistent0")

	if err := source.Open(); err != ErrWindowsCaptureUnsupported {
		t.Fatalf("expected ErrWindowsCaptureUnsupported, got %v", err)
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

func TestListInterfacesReturnsUnsupportedOnWindows(t *testing.T) {
	_, err := ListInterfaces()
	if err != ErrWindowsCaptureUnsupported {
		t.Fatalf("expected ErrWindowsCaptureUnsupported, got %v", err)
	}
}
