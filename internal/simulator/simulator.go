package simulator

import (
	"encoding/binary"
	"math/rand"
	"time"

	"github.com/DeveloperBatuhanALGUL/WardenGrid/internal/protocol/modbus"
)

type ScenarioType int

const (
	ScenarioNormalPolling ScenarioType = iota
	ScenarioUnexpectedFunctionCode
	ScenarioAbnormalWriteFrequency
	ScenarioUnauthorizedRegisterWrite
)

type Config struct {
	SourceIP        string
	DestIP          string
	UnitID          byte
	Scenario        ScenarioType
	FrameCount      int
	BaseTimestamp   int64
}

type GeneratedFrame struct {
	Raw       []byte
	SourceIP  string
	DestIP    string
	Timestamp int64
}

func Generate(cfg Config) []GeneratedFrame {
	frames := make([]GeneratedFrame, 0, cfg.FrameCount)

	switch cfg.Scenario {
	case ScenarioNormalPolling:
		frames = generateNormalPolling(cfg)
	case ScenarioUnexpectedFunctionCode:
		frames = generateUnexpectedFunctionCode(cfg)
	case ScenarioAbnormalWriteFrequency:
		frames = generateAbnormalWriteFrequency(cfg)
	case ScenarioUnauthorizedRegisterWrite:
		frames = generateUnauthorizedRegisterWrite(cfg)
	}

	return frames
}

func generateNormalPolling(cfg Config) []GeneratedFrame {
	frames := make([]GeneratedFrame, 0, cfg.FrameCount)
	for i := 0; i < cfg.FrameCount; i++ {
		txID := uint16(i + 1)
		raw := buildReadHoldingRegistersFrame(txID, cfg.UnitID, 0x00, 10)
		frames = append(frames, GeneratedFrame{
			Raw:       raw,
			SourceIP:  cfg.SourceIP,
			DestIP:    cfg.DestIP,
			Timestamp: cfg.BaseTimestamp + int64(i*2),
		})
	}
	return frames
}

func generateUnexpectedFunctionCode(cfg Config) []GeneratedFrame {
	frames := make([]GeneratedFrame, 0, cfg.FrameCount)
	for i := 0; i < cfg.FrameCount; i++ {
		txID := uint16(i + 1)
		var raw []byte
		if i == cfg.FrameCount/2 {
			raw = buildRawFunctionFrame(txID, cfg.UnitID, 0x2B)
		} else {
			raw = buildReadHoldingRegistersFrame(txID, cfg.UnitID, 0x00, 10)
		}
		frames = append(frames, GeneratedFrame{
			Raw:       raw,
			SourceIP:  cfg.SourceIP,
			DestIP:    cfg.DestIP,
			Timestamp: cfg.BaseTimestamp + int64(i*2),
		})
	}
	return frames
}

func generateAbnormalWriteFrequency(cfg Config) []GeneratedFrame {
	frames := make([]GeneratedFrame, 0, cfg.FrameCount)
	for i := 0; i < cfg.FrameCount; i++ {
		txID := uint16(i + 1)
		raw := buildWriteSingleRegisterFrame(txID, cfg.UnitID, 0x10, uint16(rand.Intn(100)))
		frames = append(frames, GeneratedFrame{
			Raw:       raw,
			SourceIP:  cfg.SourceIP,
			DestIP:    cfg.DestIP,
			Timestamp: cfg.BaseTimestamp,
		})
	}
	return frames
}

func generateUnauthorizedRegisterWrite(cfg Config) []GeneratedFrame {
	frames := make([]GeneratedFrame, 0, cfg.FrameCount)
	for i := 0; i < cfg.FrameCount; i++ {
		txID := uint16(i + 1)
		var raw []byte
		if i == cfg.FrameCount-1 {
			raw = buildWriteSingleRegisterFrame(txID, cfg.UnitID, 0xFF00, 1)
		} else {
			raw = buildReadHoldingRegistersFrame(txID, cfg.UnitID, 0x00, 10)
		}
		frames = append(frames, GeneratedFrame{
			Raw:       raw,
			SourceIP:  cfg.SourceIP,
			DestIP:    cfg.DestIP,
			Timestamp: cfg.BaseTimestamp + int64(i*2),
		})
	}
	return frames
}

func buildReadHoldingRegistersFrame(txID uint16, unitID byte, startAddr uint16, quantity uint16) []byte {
	body := make([]byte, 6)
	body[0] = byte(modbus.FuncReadHoldingRegisters)
	binary.BigEndian.PutUint16(body[1:3], startAddr)
	binary.BigEndian.PutUint16(body[3:5], quantity)
	return wrapMBAP(txID, unitID, body)
}

func buildWriteSingleRegisterFrame(txID uint16, unitID byte, addr uint16, value uint16) []byte {
	body := make([]byte, 5)
	body[0] = byte(modbus.FuncWriteSingleRegister)
	binary.BigEndian.PutUint16(body[1:3], addr)
	binary.BigEndian.PutUint16(body[3:5], value)
	return wrapMBAP(txID, unitID, body)
}

func buildRawFunctionFrame(txID uint16, unitID byte, functionCode byte) []byte {
	body := []byte{functionCode, 0x00, 0x00}
	return wrapMBAP(txID, unitID, body)
}

func wrapMBAP(txID uint16, unitID byte, body []byte) []byte {
	frame := make([]byte, 7+len(body))
	binary.BigEndian.PutUint16(frame[0:2], txID)
	binary.BigEndian.PutUint16(frame[2:4], 0)
	binary.BigEndian.PutUint16(frame[4:6], uint16(len(body)+1))
	frame[6] = unitID
	copy(frame[7:], body)
	return frame
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
