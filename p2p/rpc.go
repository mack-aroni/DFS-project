package p2p

import "net"

// RPC holds any data that is being sent between two nodes

type RPC struct {
	From    net.Addr
	Payload []byte
}
