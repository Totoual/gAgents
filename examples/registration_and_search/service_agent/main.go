package main

import (
	"fmt"
	"os"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/registration_and_search/service_agent/acts"
	"github.com/totoual/gAgents/examples/registration_and_search/service_agent/config"
	"github.com/totoual/gAgents/examples/registration_and_search/service_agent/tasks"
	"github.com/totoual/gAgents/services/kafka"
	"gopkg.in/yaml.v3"
)

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var config config.AgentConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	agent := gAgents.NewAgent(config.AgentName, config.AgentURL)
	r := acts.NewRegistrationAct(config, agent.Dispatcher)
	agent.RegisterAct(r)

	t := tasks.NewHeartbeatTask(time.Now(), 2*time.Minute, config)
	agent.TaskScheduler.AddTask(t)

	kafka, err := kafka.NewKafkaConsumerService(
		[]string{config.KafkaURL},
		agent.Dispatcher,
		acts.AgentRegisteredEventType,
	)
	if err != nil {
		// Handle error
		fmt.Println("Error creating Kafka service:", err)
		return
	}
	fmt.Println(kafka)
	proposal := acts.NewSendProposalAct(config, agent.Dispatcher, agent)
	fmt.Println(proposal)
	agent.Run()

}
