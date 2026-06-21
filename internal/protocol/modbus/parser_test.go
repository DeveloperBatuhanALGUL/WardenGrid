package modbus

import "testing"

func TestParseFrameReadHoldingRegisters(t *testing.T) {
	raw := []byte{
		0x00, 0x01,
		0x00, 0x00,
		0x00, 0x06,
		0x01,
		0x03,
		0x00, 0x6B,
		0x00, 0x03,
	}

	frame, err := ParseFrame(raw, "10.0.0.5", "10.0.0.10", 1700000000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.Header.TransactionID != 1 {
		t.Errorf("expected transaction id 1, got %d", frame.Header.TransactionID)
	}
	if frame.Header.UnitID != 1 {
		t.Errorf("expected unit id 1, got %d", frame.Header.UnitID)
	}
	if frame.Function != FuncReadHoldingRegisters {
		t.Errorf("expected function code %v, got %v", FuncReadHoldingRegisters, frame.Function)
	}
	if frame.StartAddress != 0x6B {
		t.Errorf("expected start address 0x6B, got 0x%X", frame.StartAddress)
	}
	if frame.Quantity != 3 {
		t.Errorf("expected quantity 3, got %d", frame.Quantity)
	}
	if frame.IsException {
		t.Errorf("did not expect exception frame")
	}
}

func TestParseFrameWriteSingleRegister(t *testing.T) {
	raw := []byte{
		0x00, 0x02,
		0x00, 0x00,
		0x00, 0x06,
		0x01,
		0x06,
		0x00, 0x01,
		0x00, 0x03,
	}

	frame, err := ParseFrame(raw, "10.0.0.5", "10.0.0.10", 1700000001)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.Function != FuncWriteSingleRegister {
		t.Errorf("expected function code %v, got %v", FuncWriteSingleRegister, frame.Function)
	}
	if !frame.Function.IsWrite() {
		t.Errorf("expected IsWrite to be true")
	}
	if frame.StartAddress != 1 {
		t.Errorf("expected start address 1, got %d", frame.StartAddress)
	}
	if frame.WriteValue != 3 {
		t.Errorf("expected write value 3, got %d", frame.WriteValue)
	}
}

func TestParseFrameException(t *testing.T) {
	raw := []byte{
		0x00, 0x03,
		0x00, 0x00,
		0x00, 0x03,
		0x01,
		0x83,
		0x02,
	}

	frame, err := ParseFrame(raw, "10.0.0.5", "10.0.0.10", 1700000002)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !frame.IsException {
		t.Fatalf("expected exception frame")
	}
	if frame.Exception != ExceptionIllegalDataAddress {
		t.Errorf("expected exception code %v, got %v", ExceptionIllegalDataAddress, frame.Exception)
	}
}

func TestParseFrameTooShort(t *testing.T) {
	raw := []byte{0x00, 0x01, 0x00}

	_, err := ParseFrame(raw, "10.0.0.5", "10.0.0.10", 1700000003)
	if err != ErrFrameTooShort {
		t.Fatalf("expected ErrFrameTooShort, got %v", err)
	}
}

func TestParseFrameInvalidProtocol(t *testing.T) {
	raw := []byte{
		0x00, 0x01,
		0x00, 0x01,
		0x00, 0x06,
		0x01,
		0x03,
		0x00, 0x00,
		0x00, 0x00,
	}

	_, err := ParseFrame(raw, "10.0.0.5", "10.0.0.10", 1700000004)
	if err != ErrInvalidProtocol {
		t.Fatalf("expected ErrInvalidProtocol, got %v", err)
	}
}
