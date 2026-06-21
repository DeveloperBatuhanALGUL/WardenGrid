package dnp3

type ApplicationFunctionCode byte

const (
	FuncConfirm           ApplicationFunctionCode = 0x00
	FuncRead              ApplicationFunctionCode = 0x01
	FuncWrite             ApplicationFunctionCode = 0x02
	FuncSelect            ApplicationFunctionCode = 0x03
	FuncOperate           ApplicationFunctionCode = 0x04
	FuncDirectOperate     ApplicationFunctionCode = 0x05
	FuncDirectOperateNoAck ApplicationFunctionCode = 0x06
	FuncColdRestart       ApplicationFunctionCode = 0x0D
	FuncWarmRestart       ApplicationFunctionCode = 0x0E
	FuncUnsolicitedResponse ApplicationFunctionCode = 0x82
	FuncResponse          ApplicationFunctionCode = 0x81
)

type LinkControl struct {
	Direction          bool
	Primary            bool
	FrameCountBit      bool
	FrameCountValid    bool
	FunctionCode       byte
}

type LinkHeader struct {
	Length      byte
	Control     LinkControl
	Destination uint16
	Source      uint16
	HeaderCRC   uint16
}

type TransportHeader struct {
	FIR      bool
	FIN      bool
	Sequence byte
}

type ApplicationHeader struct {
	FIR          bool
	FIN          bool
	Confirm      bool
	Unsolicited  bool
	Sequence     byte
	FunctionCode ApplicationFunctionCode
}

type Frame struct {
	Link        LinkHeader
	Transport   TransportHeader
	Application ApplicationHeader
	Payload     []byte
	SourceIP    string
	DestIP      string
	Timestamp   int64
}

func (f ApplicationFunctionCode) IsWrite() bool {
	switch f {
	case FuncWrite, FuncSelect, FuncOperate, FuncDirectOperate, FuncDirectOperateNoAck:
		return true
	default:
		return false
	}
}

func (f ApplicationFunctionCode) IsRestart() bool {
	switch f {
	case FuncColdRestart, FuncWarmRestart:
		return true
	default:
		return false
	}
}

func (f ApplicationFunctionCode) IsKnown() bool {
	switch f {
	case FuncConfirm, FuncRead, FuncWrite, FuncSelect, FuncOperate, FuncDirectOperate,
		FuncDirectOperateNoAck, FuncColdRestart, FuncWarmRestart, FuncResponse, FuncUnsolicitedResponse:
		return true
	default:
		return false
	}
}
