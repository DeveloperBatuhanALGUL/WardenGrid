//go:build windows

package capture

import (
	"errors"
	"time"
)

var ErrWindowsCaptureUnsupported = errors.New("capture: live packet capture on Windows requires Npcap and is not yet implemented")

type PcapSource struct {
	InterfaceName string
	SnapshotLen   int32
	Promiscuous   bool
	Timeout       time.Duration
}

func NewPcapSource(interfaceName string) *PcapSource {
	return &PcapSource{
		InterfaceName: interfaceName,
		SnapshotLen:   65535,
		Promiscuous:   false,
		Timeout:       time.Second,
	}
}

func (s *PcapSource) Open() error {
	return ErrWindowsCaptureUnsupported
}

func (s *PcapSource) Read() (*RawPacket, error) {
	return nil, ErrWindowsCaptureUnsupported
}

func (s *PcapSource) Close() error {
	return nil
}

func ListInterfaces() ([]InterfaceInfo, error) {
	return nil, ErrWindowsCaptureUnsupported
}
