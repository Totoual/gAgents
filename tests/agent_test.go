package tests_test

import (
	"context"
	"testing"
	"time"

	"github.com/totoual/gAgents"
)

func TestRunAgent(t *testing.T) {

	agent := gAgents.NewAgent("test", "0.0.0.0:8000")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start the agent in a goroutine so it doesn't block
	go agent.Run(ctx)

	// ... Perform your test logic here

	// Wait for the agent to stop (optional)
	<-ctx.Done()

}

func TestAgentCommunication(t *testing.T) {
	// Create two agents
	agent1 := gAgents.NewAgent("Agent1", "localhost:8000")
	agent2 := gAgents.NewAgent("Agent2", "localhost:8001")

	// Start the agents
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go agent1.Run(ctx)
	go agent2.Run(ctx)

	// Start goroutines to consume messages
	go agent1.ConsumeInMessages()
	go agent2.ConsumeInMessages()

	// Send a message from Agent1 to Agent2
	response, err := agent1.DoSendMessage("localhost:8001", "Agent2", "Hello from Agent1!")
	if err != nil {
		t.Fatalf("Agent1 failed to send message: %v", err)
	}

	// Check if the response status is "OK"
	if response.Status != "OK" {
		t.Fatalf("Unexpected response status: %s", response.Status)
	}

	// Wait for a short time to allow for message processing
	time.Sleep(100 * time.Millisecond)
}
