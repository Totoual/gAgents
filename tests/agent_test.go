package tests_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/totoual/gAgents"
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

func TestRunAgent(t *testing.T) {

	agent := gAgents.NewAgent("test", "0.0.0.0:8000")
	ctx, cancel := context.WithCancel(context.Background())

	// Start the agent in a goroutine so it doesn't block
	go agent.Run(ctx)

	// ... Perform your test logic here

	// Wait for the agent to stop (optional)
	cancel()

}

func TestAgentCommunication(t *testing.T) {
	// Create two agents
	agent1 := gAgents.NewAgent("Agent1", "0.0.0.0:8001")
	agent2 := gAgents.NewAgent("Agent2", "0.0.0.0:8002")

	// Define a custom handler for greeting messages
	greetingHandler := &GreetingHandler{}

	// Register the handler for the "greet" message type
	agent2.RegisterHandler("greet", greetingHandler)

	// Start Agent1 in a goroutine
	ctx1, cancel1 := context.WithCancel(context.Background())
	go agent1.Run(ctx1)

	// Start Agent2 in a goroutine
	ctx2, cancel2 := context.WithCancel(context.Background())
	go agent2.Run(ctx2)

	// Allow time for the servers to start
	time.Sleep(time.Second)

	// Simulate sending a message from Agent1 to Agent2
	agent1.OutMessageQueue <- gAgents.Message{
		Receiver: agent2.Addr,
		Type:     "greet",
		Content:  []byte("Hello!"),
	}

	// Allow time for message processing
	time.Sleep(time.Second)

	// Add any assertions or checks here based on the expected behavior

	// Cleanup and stop the agents
	cancel1()
	cancel2()
}
