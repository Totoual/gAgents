# gAgents: Micro Agent Framework

gAgents is a micro-framework that provides a foundation for building task-oriented agents in Go.
To get started, follow these simple steps:

## Installation

Get started by installing and configuring your Go environment. The current Go version is `1.20`.

1. Open your terminal and navigate to your Go project.
2. Use the `go get` command to install the framework:

```bash
go get github.com/totoual/gAgents
```

This will download the main branch of the framework and add it to your project's dependencies.

## Getting Started

Now that you have the framework installed, you can start building agents. Here's a basic example to help you get going:

### Creating a Message Struct

Firstly, you need to create a message struct. Create a file called `protocol.go`:

```go
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
	Protocol string `json:"protocol"`
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

func (t TestMessage) GetPerformative() int{
    return 0
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
```

### Setting up the Agent

In your `agent.go` file:

```go
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
		Protocol: "greet",
		Greeting: "Hello!",
	}

	// Send a message
	agent.OutMessageQueue <- *gAgents.NewEnvelope(message)

	// Allow time for message processing
	time.Sleep(time.Second)

	agent.Cancel()
}
```

Run this agent:

```bash
go run main.go
```

## Messaging System

Learn more about the messaging system in [message.md](./docs/message.md).

## Task Scheduler

Learn more about the tasks and the scheduler in [task.md](./docs//task.md).
