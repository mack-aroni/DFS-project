package p2p

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// RPC holds any data that is being sent between two nodes
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
