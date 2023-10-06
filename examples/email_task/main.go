package main

import (
	"fmt"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/email_task/task"
)

func main() {

	agent := gAgents.NewAgent("Alice", "localhost:8000")
	fmt.Printf(agent.Addr)
	emailService := task.NewEmailService(agent.TaskScheduler)
	agent.RegisterService(emailService)
	agent.Run()

}
