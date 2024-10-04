package p2p

import "net"

// Peer represents a remote node
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// Transport holds communication between nodes
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
