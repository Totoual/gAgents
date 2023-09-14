package main

import (
	"fmt"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
)

// CustomAct is an example of a custom act.
type CustomAct struct {
	count int
}

// Perform is the method called when the act is performed.
func (a *CustomAct) Perform(agent *gAgents.Agent) {
	// Perform action...
	a.count++
	fmt.Printf("CustomAct performed. Count: %d\n", a.count)
}

// GetInterval returns the interval at which the act should be performed.
func (a *CustomAct) GetInterval() time.Duration {
	return time.Second * 2 // Perform every 2 seconds
}

func main() {
	// Create a new agent
	agent := gAgents.NewAgent("Alice", "localhost:8000")

	// Create a custom act
	myAct := &CustomAct{}

	// Register the act with the agent
	agent.RegisterAct(myAct)

	// Start performing acts
	go agent.PerformActs()

	// Let the agent run for a while to observe act execution
	time.Sleep(time.Second * 10)

	// Stop the agent (Optional)
	agent.Cancel()
}
