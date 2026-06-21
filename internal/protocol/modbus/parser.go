package modbus

import (
	"encoding/binary"
	"errors"
)

var (
	ErrFrameTooShort   = errors.New("modbus: frame shorter than mbap header")
	ErrPayloadTooShort = errors.New("modbus: payload shorter than declared length")
	ErrInvalidProtocol = errors.New("modbus: protocol identifier not zero")
)

const mbapHeaderSize = 7

func ParseFrame(raw []byte, sourceIP string, destIP string, timestamp int64) (*Frame, error) {
	if len(raw) < mbapHeaderSize {
		return nil, ErrFrameTooShort
	}

	header := MBAPHeader{
		TransactionID: binary.BigEndian.Uint16(raw[0:2]),
		ProtocolID:    binary.BigEndian.Uint16(raw[2:4]),
		Length:        binary.BigEndian.Uint16(raw[4:6]),
		UnitID:        raw[6],
	}

	if header.ProtocolID != 0 {
		return nil, ErrInvalidProtocol
	}

	expectedTotal := mbapHeaderSize + int(header.Length) - 1
	if len(raw) < expectedTotal {
		return nil, ErrPayloadTooShort
	}

	body := raw[mbapHeaderSize:expectedTotal]
	if len(body) < 1 {
		return nil, ErrFrameTooShort
	}

	frame := &Frame{
		Header:    header,
		SourceIP:  sourceIP,
		DestIP:    destIP,
		Timestamp: timestamp,
	}

	rawFunc := body[0]
	if rawFunc&0x80 != 0 {
		frame.IsException = true
		frame.Function = FunctionCode(rawFunc & 0x7F)
		if len(body) >= 2 {
			frame.Exception = ExceptionCode(body[1])
		}
		return frame, nil
	}

	frame.Function = FunctionCode(rawFunc)
	frame.RawPayload = body[1:]

	switch frame.Function {
	case FuncReadCoils, FuncReadDiscreteInputs, FuncReadHoldingRegisters, FuncReadInputRegisters:
		if len(frame.RawPayload) >= 4 {
			frame.StartAddress = binary.BigEndian.Uint16(frame.RawPayload[0:2])
			frame.Quantity = binary.BigEndian.Uint16(frame.RawPayload[2:4])
		}
	case FuncWriteSingleCoil, FuncWriteSingleRegister:
		if len(frame.RawPayload) >= 4 {
			frame.StartAddress = binary.BigEndian.Uint16(frame.RawPayload[0:2])
			frame.WriteValue = binary.BigEndian.Uint16(frame.RawPayload[2:4])
		}
	case FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		if len(frame.RawPayload) >= 4 {
			frame.StartAddress = binary.BigEndian.Uint16(frame.RawPayload[0:2])
			frame.Quantity = binary.BigEndian.Uint16(frame.RawPayload[2:4])
		}
	}

	return frame, nil
}
