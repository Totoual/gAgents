package main

import (
	"log"

	gAgents "github.com/totoual/gAgents/agent"
)

type GreetingHandler struct {
	// You can include any necessary properties for the handler here
}

func (h *GreetingHandler) HandleMessage(message gAgents.Message) {
	// Implement the handling logic for greeting messages here
	if message.Type == "greet" {
		// Assuming the greeting message format is "Hello, {Receiver}!"
		log.Printf("Received greeting: %s\n", string(message.Content))

		// Send a response back to the sender
	} else {
		log.Printf("Unsupported message type: %s\n", message.Type)
	}
}

func main() {

	a := gAgents.NewAgent("Test", "0.0.0.0:8002")
	a.RegisterHandler("greet", &GreetingHandler{})
	a.Run()

}
