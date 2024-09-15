package p2p

// Peer represents a remote node
type Peer interface {
	Close() error
}

// Transport holds communication between nodes
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
