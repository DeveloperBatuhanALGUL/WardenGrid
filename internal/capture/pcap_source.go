//go:build !windows

package capture

import (
	"errors"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var ErrSourceNotOpen = errors.New("capture: pcap source not open")

type PcapSource struct {
	InterfaceName string
	SnapshotLen   int32
	Promiscuous   bool
	Timeout       time.Duration
	handle        *pcap.Handle
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
	handle, err := pcap.OpenLive(s.InterfaceName, s.SnapshotLen, s.Promiscuous, s.Timeout)
	if err != nil {
		return err
	}
	s.handle = handle
	return nil
}

func (s *PcapSource) Read() (*RawPacket, error) {
	if s.handle == nil {
		return nil, ErrSourceNotOpen
	}

	data, captureInfo, err := s.handle.ReadPacketData()
	if err != nil {
		return nil, err
	}

	packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
	sourceIP, destIP := extractIPs(packet)
	payload := extractTransportPayload(packet)

	return &RawPacket{
		Data:      payload,
		SourceIP:  sourceIP,
		DestIP:    destIP,
		Timestamp: captureInfo.Timestamp.Unix(),
	}, nil
}

func (s *PcapSource) Close() error {
	if s.handle != nil {
		s.handle.Close()
		s.handle = nil
	}
	return nil
}

func extractIPs(packet gopacket.Packet) (string, string) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return "", ""
	}

	ip, ok := ipLayer.(*layers.IPv4)
	if !ok {
		return "", ""
	}

	return ip.SrcIP.String(), ip.DstIP.String()
}

func extractTransportPayload(packet gopacket.Packet) []byte {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return nil
	}

	tcp, ok := tcpLayer.(*layers.TCP)
	if !ok {
		return nil
	}

	return tcp.Payload
}
