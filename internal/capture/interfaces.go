//go:build !windows

package capture

import "github.com/google/gopacket/pcap"

func ListInterfaces() ([]InterfaceInfo, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	result := make([]InterfaceInfo, 0, len(devices))
	for _, dev := range devices {
		addresses := make([]string, 0, len(dev.Addresses))
		for _, addr := range dev.Addresses {
			if addr.IP != nil {
				addresses = append(addresses, addr.IP.String())
			}
		}
		result = append(result, InterfaceInfo{
			Name:        dev.Name,
			Description: dev.Description,
			Addresses:   addresses,
		})
	}

	return result, nil
}
