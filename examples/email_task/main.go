package main

import (
	"fmt"

	gAgents "github.com/totoual/gAgents/agent"
)

func main() {

	agent := gAgents.NewAgent("Alice", "localhost:8000")
	fmt.Printf(agent.Addr)

}
