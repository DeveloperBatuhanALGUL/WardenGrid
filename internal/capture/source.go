package capture

type RawPacket struct {
	Data      []byte
	SourceIP  string
	DestIP    string
	Timestamp int64
}

type Source interface {
	Open() error
	Read() (*RawPacket, error)
	Close() error
}
