package main

import (
	"fmt"
	"log"

	"github.com/mack-aroni/DFS-project/p2p"
)

func OnPeer(peer p2p.Peer) error {
	fmt.Println("doing some logic with peer outside of TCPTransport")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		Handshake:  p2p.NOPHandshakeFunc,
		Decoder:    p2p.DefaultDecoder{},
		OnPeer:     OnPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
