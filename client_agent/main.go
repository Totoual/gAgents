package main

import (
	"log"

	"github.com/totoual/gAgents"
)

func main() {

	a := gAgents.NewAgent("test2", "0.0.0.0:8001")
	res, err := a.DoSendMessage("localhost:8000", "Agent1", "This is a test!")
	if err != nil {
		log.Fatalf("Could not send message: %v\n", err)
	}

	log.Printf("The response is: %v\n", res)

}
