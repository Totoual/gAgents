package main

import (
	"log"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/server_client_communication/protocol"
)

func main() {
	a := gAgents.NewAgent("test2", "0.0.0.0:8003")

	// Start the agent in a goroutine
	go a.Run()

	time.Sleep(1 * time.Second)

	message := protocol.TestMessage{
		Receiver: "0.0.0.0:8002",
		Sender:   "0.0.0.0:8003",
		Type:     "greet",
		Greeting: "Hello!",
	}

	// Send a message
	a.OutMessageQueue <- *gAgents.NewEnvelope(message)
	log.Printf("Added the message in the Queue")

	// Wait for the wait group to signal that the agent.Run() is done
	a.Cancel()

}
