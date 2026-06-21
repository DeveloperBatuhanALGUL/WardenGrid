package dnp3

import "testing"

func TestParseFrameReadRequest(t *testing.T) {
	raw := []byte{
		0x05, 0x64,
		0x0B,
		0xC4,
		0x46, 0x00,
		0x40, 0x00,
		0xA3, 0xFE,
		0xC0,
		0xCD,
		0x01,
		0x3C, 0x01,
	}

	frame, err := ParseFrame(raw, "10.0.0.5", "10.0.0.10", 1700000000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.Link.Destination != 0x0046 {
		t.Errorf("expected destination 0x0046, got 0x%X", frame.Link.Destination)
	}
	if frame.Link.Source != 0x0040 {
		t.Errorf("expected source 0x0040, got 0x%X", frame.Link.Source)
	}
	if !frame.Link.Control.Primary {
		t.Errorf("expected primary bit set")
	}
	if !frame.Transport.FIR || !frame.Transport.FIN {
		t.Errorf("expected FIR and FIN set on single-frame transport header")
	}
	if frame.Application.FunctionCode != FuncRead {
		t.Errorf("expected function code FuncRead, got 0x%X", byte(frame.Application.FunctionCode))
	}
	if !frame.Application.FunctionCode.IsKnown() {
		t.Errorf("expected known function code")
	}
}

func TestParseFrameInvalidStartBytes(t *testing.T) {
	raw := []byte{
		0x00, 0x00,
		0x0B,
		0xC4,
		0x46, 0x00,
		0x40, 0x00,
		0xA3, 0xFE,
	}

	_, err := ParseFrame(raw, "10.0.0.5", "10.0.0.10", 1700000000)
	if err != ErrInvalidStartBytes {
		t.Fatalf("expected ErrInvalidStartBytes, got %v", err)
	}
}

func TestParseFrameTooShort(t *testing.T) {
	raw := []byte{0x05, 0x64, 0x0B}

	_, err := ParseFrame(raw, "10.0.0.5", "10.0.0.10", 1700000000)
	if err != ErrFrameTooShort {
		t.Fatalf("expected ErrFrameTooShort, got %v", err)
	}
}

func TestParseFrameWriteFunctionIsWrite(t *testing.T) {
	raw := []byte{
		0x05, 0x64,
		0x0B,
		0xC4,
		0x46, 0x00,
		0x40, 0x00,
		0xA3, 0xFE,
		0xC0,
		0xCD,
		0x05,
		0x0C, 0x01,
	}

	frame, err := ParseFrame(raw, "10.0.0.5", "10.0.0.10", 1700000000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.Application.FunctionCode != FuncDirectOperate {
		t.Errorf("expected FuncDirectOperate, got 0x%X", byte(frame.Application.FunctionCode))
	}
	if !frame.Application.FunctionCode.IsWrite() {
		t.Errorf("expected IsWrite to be true")
	}
}
