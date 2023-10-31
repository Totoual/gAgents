package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/registry/services/registry"
	healthcheck "github.com/totoual/gAgents/registry/tasks/health_check"
	chatgpt "github.com/totoual/gAgents/services/chatGPT"
	"github.com/totoual/gAgents/services/kafka"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AgentName string   `yaml:"agent_name"`
	AgentURL  string   `yaml:"agent_url"`
	KafkaURL  string   `yaml:"kafka_url"`
	Topics    []string `yaml:"topics"`
}

func init() {
	err := godotenv.Load() // Load .env file from the current directory
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
	reg := registry.NewRegistryService(agent.Dispatcher, agent.TaskScheduler)
	heartbeat := healthcheck.NewHeartbeatTask(time.Now(), 30*time.Minute, reg)
	gpt := chatgpt.NewGPTClient(
		agent.Dispatcher, "gpt-4", registry.SearchEvent, registry.TopicSuggestionEvent, config.Topics,
	)
	fmt.Printf(gpt.ApiKey)

	agent.TaskScheduler.AddTask(heartbeat)
	k, err := kafka.NewKafkaProducerService(
		[]string{config.KafkaURL},
		agent.Dispatcher,
		config.Topics,
	)
	if err != nil {
		// Handle error
		fmt.Println("Error creating Kafka service:", err)
		return
	}
	fmt.Println(k)
	agent.RegisterService(reg)
	fmt.Printf(agent.Addr)
	agent.Run()

}
