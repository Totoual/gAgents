package main

import (
	"os"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/registration_and_search/search_agent/config"
	userinput "github.com/totoual/gAgents/examples/registration_and_search/search_agent/services/user_input"
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
	user_input := userinput.NewUserInteractionService(
		agent.Dispatcher,
		config,
	)
	agent.RegisterService(user_input)
	agent.Run()

}
