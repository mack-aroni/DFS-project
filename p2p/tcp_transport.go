package p2p

import (
	"errors"
	"fmt"
	"net"
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

// Send implements the Peer interface
func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

// RemoteAddr implements the Peer interface and returns
// the remote address of its underlying connection
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
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
}

// TCPTransport constructor
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan RPC),
	}
}

// Consume implements Transport interface, which returns
// a read-only channel for messages from another Peer
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

// Close implements Transport interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implements Transport interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

// TCPTransport listens over TCP network for connections
func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	fmt.Printf("TCP transport listening on port: %s\n", t.ListenAddr)

	return nil
}

// Loops the listener for connection events
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		go t.handleConn(conn, false)
	}
}

// temporary
type Temp struct {
}

// Directly handles creating new TCPPeer from a connection
func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

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
