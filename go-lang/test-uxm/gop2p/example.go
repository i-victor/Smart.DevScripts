package main

import (
	"fmt"
	"time"

	"github.com/lucasmenendez/gop2p"
)

func main() {
	// Creating main node with debug mode equal to false. Then set individual
	// handlers.
	p2p := gop2p.InitNode(5001, false)
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

	// Creating peer on localhost 5002 port.
	go func() {
		// Wait for main node initialization.
		time.Sleep(time.Second)
		// Get main peer and create node in debug mode. To create an entry peer
		// manually, use CreatePeer function.
		entry := p2p.Self
		node := gop2p.InitNode(5002, true)

		// Connect to main node peer.
		node.Connect(entry)
		// Wait and broadcast message.
		time.Sleep(time.Second)
		node.Broadcast([]byte("Hello peers, from 2nd"))
		// Wait and disconnect
		time.Sleep(2 * time.Second)
		node.Disconnect()
	}()

	// Create peer on localhost 5003 port.
	go func() {
		time.Sleep(time.Second)
		entry := p2p.Self

		node := gop2p.InitNode(5003, false)

		node.Connect(entry)
		time.Sleep(2 * time.Second)
		node.Disconnect()
	}()

	// Wait and broadcast. Broadcast fail is expected.
	time.Sleep(6 * time.Second)
	p2p.Broadcast([]byte("Hello peers, from 1st"))
	// Wait and disconnect.
	time.Sleep(2 * time.Second)
	p2p.Disconnect()

	// Output:[main handler] -> Connected
	//[main handler] -> Connected
	//[main handler] -> Message: Hello peers!
	//[main handler] -> Disconnected
	//[main handler] -> Disconnected
}
