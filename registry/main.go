package main

import (
	"fmt"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/registry/services/kafka"
	"github.com/totoual/gAgents/registry/services/registry"
)

func main() {

	agent := gAgents.NewAgent("Alice", "localhost:8010")
	registry := registry.NewRegistryService(agent.Dispatcher)
	kafka, err := kafka.NewKafkaService([]string{"localhost:9092"}, agent.Dispatcher)
	if err != nil {
		// Handle error
		fmt.Println("Error creating Kafka service:", err)
		return
	}
	fmt.Println(kafka)
	agent.RegisterService(registry)
	fmt.Printf(agent.Addr)
	agent.Run()

}
