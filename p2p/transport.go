package p2p

// Peer represents a remote node
type Peer interface {
	Close() error
}

// Transport holds communication between nodes
type Transport interface {
	ListenAndAccept() error
	consume() <-chan RPC
}
