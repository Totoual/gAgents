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

	// Start Agent1 and Agent2 in separate goroutines
	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())

	go agent1.Run(ctx1)
	go agent2.Run(ctx2)

	// Give some time for the servers to start
	time.Sleep(100 * time.Millisecond)

	// Agent1 sends a message to Agent2
	response, err := agent1.DoSendMessage(agent2.Addr, "Agent2", "Hello from Agent1!")
	if err != nil {
		t.Errorf("Agent1 failed to send message: %v", err)
	}
	t.Logf(response.Status)

	// Agent2 should receive the message
	select {
	case msg := <-agent2.InMessageQueue:
		if msg.Content != "Hello from Agent1!" {
			t.Errorf("Agent2 received unexpected message content: %s", msg.Content)
		}
	default:
		t.Error("Agent2 did not receive any message")
	}

	// Agent2 sends a message to Agent1
	agent2.SendAsyncMessage("Agent1", "Hello from Agent2!")

	// Agent1 should receive the message
	select {
	case msg := <-agent1.InMessageQueue:
		if msg.Content != "Hello from Agent2!" {
			t.Errorf("Agent1 received unexpected message content: %s", msg.Content)
		}
	default:
		t.Error("Agent1 did not receive any message")
	}

	// Clean up
	cancel1()
	cancel2()
}
