# gAgents: Micro Agent Framework

gAgents is a micro-framework that provides a foundation for building agents in Go.
To get started, follow this simple setps:

### Installation
Get started by installing and configuring your Go env. Current go version is `1.20`

1. Open your terminal and navigate to your Go project
2. Use the `go get` command to install the framework:

```bash
go get github.com/totoual/gAgents/agent
```

This will download the main branch of the framework and add it to your project's
dependencies.

### Getting Started

Now that you have the framework installed, you can start building agents. Here's a really basic example to help you get going:

```go
package main

import(
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


func main(){
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
```
Run this agent:

```bash
go run main.go
```