package main

import (
	"fmt"
	"os"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/registry/services/registry"
	healthcheck "github.com/totoual/gAgents/registry/tasks/health_check"
	"github.com/totoual/gAgents/services/kafka"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AgentName string   `yaml:"agent_name"`
	AgentURL  string   `yaml:"agent_url"`
	KafkaURL  string   `yaml:"kafka_url"`
	Topics    []string `yaml:"topics"`
}

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	agent := gAgents.NewAgent(config.AgentName, config.AgentURL)
	registry := registry.NewRegistryService(agent.Dispatcher, agent.TaskScheduler)
	heartbeat := healthcheck.NewHeartbeatTask(time.Now(), 30*time.Minute, registry)
	agent.TaskScheduler.AddTask(heartbeat)
	kafka, err := kafka.NewKafkaProducerService(
		[]string{config.KafkaURL},
		agent.Dispatcher,
		config.Topics,
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
