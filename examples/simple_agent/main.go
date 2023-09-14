package main

import (
	"time"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/simple_agent/protocol"
)

func main() {
	// Create a new agent
	agent := gAgents.NewAgent("Alice", "localhost:8000")

	// Here you can register acts, handlers and tasks
	agent.RegisterHandler("greet", &protocol.GreetingHandler{})

	// Run the agent.
	go agent.Run() // Need to run the agent in a different routine so we can send a message for the example:

	message := protocol.TestMessage{
		Receiver: "0.0.0.0:8002",
		Sender:   "0.0.0.0:8003",
		Type:     "greet",
		Greeting: "Hello!",
	}

	// Send a message
	agent.OutMessageQueue <- *gAgents.NewEnvelope(message)

	// Allow time for message processing
	time.Sleep(time.Second)

	agent.Cancel()
}
