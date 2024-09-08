package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	conn net.Conn
	// dialing connections are outbound, accepting connection are not
	outbound bool
}

// TCPPeer constructor
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// Options struct for TCPTransport
type TCPTransportOpts struct {
	ListenAddr string
	Handshake  HandshakeFunc
	Decoder    Decoder
	OnPeer     func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcCh    chan RPC

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

// TCPTransport constructor
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan RPC),
	}
}

// Consume implements Transport interface,
// returns a read-only channel for messages from another Peer
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

// TCPTransport listens over TCP network for connections
func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

// Loops the listener for connection events
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		fmt.Printf("Handling new incoming connection: %+v\n", conn)

		go t.handleConn(conn)
	}
}

// temporary
type Temp struct {
}

// Directly handles creating new TCPPeer from a connection
func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err := t.Handshake(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)

		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr()

		t.rpcCh <- rpc
	}

}
