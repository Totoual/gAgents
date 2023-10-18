package main

import (
	"fmt"
	"os"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/registry/services/registry"
	"github.com/totoual/gAgents/services/kafka"
	"gopkg.in/yaml.v3"
)

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var config AgentConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	agent := gAgents.NewAgent(config.AgentName, config.AgentURL)
	registry := registry.NewRegistryService(agent.Dispatcher)
	kafka, err := kafka.NewKafkaConsumerService(
		[]string{config.KafkaURL},
		agent.Dispatcher,
	)
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
