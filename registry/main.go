package main

import (
	"fmt"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/registry/services/registry"
)

func main() {

	agent := gAgents.NewAgent("Alice", "localhost:8000")
	registry := registry.NewRegistryService(agent.Dispatcher)
	agent.RegisterService(registry)
	fmt.Printf(agent.Addr)
	agent.Run()

}
