package main

import (
	"github.com/totoual/gAgents/agent" // Importing your custom package
)

func main() {
	agent1 := agent.NewAgent("Agent1")
	agent2 := agent.NewAgent("Agent2")

	go agent1.Run()
	go agent2.Run()

	message := agent.Message{
		Sender:  "Agent1",
		Content: "Hello from Agent1!",
	}

	agent1.SendMessage(message, agent2)

	// agent1.Exit()
	// agent2.Exit()
}
