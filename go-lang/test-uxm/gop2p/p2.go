package main

import (
	"fmt"
	"time"

	"github.com/lucasmenendez/gop2p"
)

func main() {
	// Creating main node with debug mode equal to false. Then set individual
	// handlers.
	p2p := gop2p.InitNode(5002, false)
	// Wait for connections.

	// Set a connection handler
	p2p.OnConnection(func(_ gop2p.Peer) {
		fmt.Printf("[main handler] -> Connected\n")
	})

	// Set a message handler.
	p2p.OnMessage(func(msg []byte, _ gop2p.Peer) {
		fmt.Printf("[main handler] -> Message: %s\n", string(msg))
	})

	// Set a disconnection handler
	p2p.OnDisconnection(func(_ gop2p.Peer) {
		fmt.Printf("[main handler] -> Disconnected\n")
	})

	defer p2p.Disconnect()

	for {
		// Wait and broadcast. Broadcast fail is expected.
		time.Sleep(10 * time.Second)
		p2p.Broadcast([]byte("Hello peers, from 1st"))
	}

}
