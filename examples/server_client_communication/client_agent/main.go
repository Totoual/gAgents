package main

import (
	"log"
	"sync"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
)

func main() {
	a := gAgents.NewAgent("test2", "0.0.0.0:8003")

	// Create a wait group
	var wg sync.WaitGroup

	// Add 1 to the wait group
	wg.Add(1)

	// Start the agent in a goroutine
	go func() {
		defer wg.Done() // When the agent.Run() exits, signal that it's done
		a.Run()
	}()

	time.Sleep(1 * time.Second)

	// Send a message
	a.OutMessageQueue <- gAgents.Message{
		Receiver: "0.0.0.0:8002",
		Sender:   "0.0.0.0:8003",
		Type:     "greet",
		Content:  []byte("Hello!"),
	}
	log.Printf("Added the message in the Queue")

	// Wait for the wait group to signal that the agent.Run() is done
	wg.Wait()

}
