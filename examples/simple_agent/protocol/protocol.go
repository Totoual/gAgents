package protocol

import (
	"encoding/json"
	"fmt"
	"log"

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
	message, ok := m.(*TestMessage)
	if !ok {
		fmt.Errorf("expected a CustomMessage")
	}
	if err != nil {
		fmt.Println("Error deserializing message:", err)
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
