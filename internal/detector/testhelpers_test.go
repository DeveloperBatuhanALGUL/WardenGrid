package detector

import "github.com/DeveloperBatuhanALGUL/WardenGrid/internal/protocol/modbus"

func mockFrame(sourceIP string, timestamp int64) modbus.Frame {
	return modbus.Frame{
		Function:     modbus.FuncReadHoldingRegisters,
		SourceIP:     sourceIP,
		DestIP:       "10.0.0.10",
		Timestamp:    timestamp,
		StartAddress: 0,
		Quantity:     1,
	}
}
