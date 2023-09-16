package tests_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
)

type TestMessage struct {
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
	Protocol string `json:"type"`
	Greeting string `json:"greeting"`
}

func (t TestMessage) GetReceiver() string {
	return t.Receiver
}
func (t TestMessage) GetSender() string {
	return t.Sender
}
func (t TestMessage) GetProtocol() string {
	return t.Protocol
}

func (t TestMessage) Serialize() ([]byte, error) {
	return json.Marshal(t)
}

type GreetingHandler struct {
	// You can include any necessary properties for the handler here
}

func (h *GreetingHandler) HandleMessage(envelope gAgents.Envelope) {
	m, err := envelope.ToMessage(&TestMessage{})
	if err != nil {
		fmt.Println("Error deserializing message:", err)
	}
	message, ok := m.(*TestMessage)
	if !ok {
		fmt.Errorf("expected a CustomMessage")
	}

	// Implement the handling logic for greeting messages here
	if message.Protocol == "greet" {
		// Assuming the greeting message format is "Hello, {Receiver}!"
		log.Printf("Received greeting: %s\n", message.Greeting)

		// Send a response back to the sender
	} else {
		log.Printf("Unsupported message type: %s\n", message.Protocol)
	}
}

func TestRunAgent(t *testing.T) {

	agent := gAgents.NewAgent("test", "0.0.0.0:8000")

	// Start the agent in a goroutine so it doesn't block
	go agent.Run()

	// ... Perform your test logic here

	// Wait for the agent to stop (optional)
	agent.Cancel()

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
	go agent1.Run()

	// Start Agent2 in a goroutine
	go agent2.Run()

	// Allow time for the servers to start
	time.Sleep(time.Second)

	message := TestMessage{
		Receiver: agent2.Addr,
		Sender:   agent1.Addr,
		Protocol: "greet",
		Greeting: "Hello!",
	}
	// Simulate sending a message from Agent1 to Agent2
	agent1.OutMessageQueue <- *gAgents.NewEnvelope(message)

	// Allow time for message processing
	time.Sleep(time.Second)

	// Add any assertions or checks here based on the expected behavior

	// Cleanup and stop the agents
	agent1.Cancel()
	agent2.Cancel()
}
