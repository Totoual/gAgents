package main

import (
	"log"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
)

type GreetingHandler struct {
	// You can include any necessary properties for the handler here
}

func (h *GreetingHandler) HandleMessage(message gAgents.Message) {
	// Implement the handling logic for greeting messages here
	if message.Type == "greet" {
		// Assuming the greeting message format is "Hello, {Receiver}!"
		log.Printf("Received greeting: %s\n", message.Content)

		// Send a response back to the sender
	} else {
		log.Printf("Unsupported message type: %s\n", message.Type)
	}
}

func main() {
	// Create a new agent
	agent := gAgents.NewAgent("Alice", "localhost:8000")

	greetingHandler := &GreetingHandler{}
	// Here you can register acts, handlers and tasks
	agent.RegisterHandler("greet", greetingHandler)

	// Run the agent.
	go agent.Run() // Need to run the agent in a different routine so we can send a message for the example:

	agent.OutMessageQueue <- gAgents.Message{
		Receiver: agent.Addr,
		Type:     "greet",
		Content:  []byte("Hello!"),
	}

	// Allow time for message processing
	time.Sleep(time.Second)

	agent.Cancel()
}
