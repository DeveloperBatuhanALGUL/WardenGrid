package dnp3

import (
	"encoding/binary"
	"errors"
)

var (
	ErrInvalidStartBytes = errors.New("dnp3: invalid start bytes, expected 0x0564")
	ErrFrameTooShort      = errors.New("dnp3: frame shorter than link header")
	ErrNoTransportHeader  = errors.New("dnp3: missing transport header")
	ErrNoApplicationHeader = errors.New("dnp3: missing application header")
)

const linkHeaderSize = 10

func ParseFrame(raw []byte, sourceIP string, destIP string, timestamp int64) (*Frame, error) {
	if len(raw) < linkHeaderSize {
		return nil, ErrFrameTooShort
	}

	if raw[0] != 0x05 || raw[1] != 0x64 {
		return nil, ErrInvalidStartBytes
	}

	link := LinkHeader{
		Length:      raw[2],
		Control:     parseLinkControl(raw[3]),
		Destination: binary.LittleEndian.Uint16(raw[4:6]),
		Source:      binary.LittleEndian.Uint16(raw[6:8]),
		HeaderCRC:   binary.LittleEndian.Uint16(raw[8:10]),
	}

	frame := &Frame{
		Link:      link,
		SourceIP:  sourceIP,
		DestIP:    destIP,
		Timestamp: timestamp,
	}

	body := raw[linkHeaderSize:]
	if len(body) < 1 {
		return frame, nil
	}

	frame.Transport = parseTransportHeader(body[0])
	body = body[1:]

	if len(body) < 1 {
		return frame, nil
	}

	frame.Application = parseApplicationHeader(body)
	if len(body) > 2 {
		frame.Payload = body[2:]
	}

	return frame, nil
}

func parseLinkControl(b byte) LinkControl {
	return LinkControl{
		Direction:       b&0x80 != 0,
		Primary:         b&0x40 != 0,
		FrameCountBit:   b&0x20 != 0,
		FrameCountValid: b&0x10 != 0,
		FunctionCode:    b & 0x0F,
	}
}

func parseTransportHeader(b byte) TransportHeader {
	return TransportHeader{
		FIR:      b&0x40 != 0,
		FIN:      b&0x80 != 0,
		Sequence: b & 0x3F,
	}
}

func parseApplicationHeader(body []byte) ApplicationHeader {
	control := body[0]
	header := ApplicationHeader{
		FIR:      control&0x80 != 0,
		FIN:      control&0x40 != 0,
		Confirm:  control&0x20 != 0,
		Unsolicited: control&0x10 != 0,
		Sequence: control & 0x0F,
	}

	if len(body) >= 2 {
		header.FunctionCode = ApplicationFunctionCode(body[1])
	}

	return header
}
