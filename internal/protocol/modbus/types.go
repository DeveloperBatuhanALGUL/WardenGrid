package modbus

type FunctionCode byte

const (
	FuncReadCoils              FunctionCode = 0x01
	FuncReadDiscreteInputs     FunctionCode = 0x02
	FuncReadHoldingRegisters   FunctionCode = 0x03
	FuncReadInputRegisters     FunctionCode = 0x04
	FuncWriteSingleCoil        FunctionCode = 0x05
	FuncWriteSingleRegister    FunctionCode = 0x06
	FuncWriteMultipleCoils     FunctionCode = 0x0F
	FuncWriteMultipleRegisters FunctionCode = 0x10
)

type ExceptionCode byte

const (
	ExceptionIllegalFunction        ExceptionCode = 0x01
	ExceptionIllegalDataAddress     ExceptionCode = 0x02
	ExceptionIllegalDataValue       ExceptionCode = 0x03
	ExceptionServerDeviceFailure    ExceptionCode = 0x04
)

type MBAPHeader struct {
	TransactionID uint16
	ProtocolID    uint16
	Length        uint16
	UnitID        byte
}

type Frame struct {
	Header       MBAPHeader
	Function     FunctionCode
	IsException  bool
	Exception    ExceptionCode
	StartAddress uint16
	Quantity     uint16
	WriteValue   uint16
	RawPayload   []byte
	SourceIP     string
	DestIP       string
	Timestamp    int64
}

func (f FunctionCode) IsWrite() bool {
	switch f {
	case FuncWriteSingleCoil, FuncWriteSingleRegister, FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		return true
	default:
		return false
	}
}

func (f FunctionCode) IsRead() bool {
	switch f {
	case FuncReadCoils, FuncReadDiscreteInputs, FuncReadHoldingRegisters, FuncReadInputRegisters:
		return true
	default:
		return false
	}
}

func (f FunctionCode) IsKnown() bool {
	switch f {
	case FuncReadCoils, FuncReadDiscreteInputs, FuncReadHoldingRegisters, FuncReadInputRegisters,
		FuncWriteSingleCoil, FuncWriteSingleRegister, FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		return true
	default:
		return false
	}
}
