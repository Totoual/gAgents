package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/totoual/gAgents"
)

func main() {
	a := gAgents.NewAgent("test2", "0.0.0.0:8003")
	ctx := context.Background()

	// Create a wait group
	var wg sync.WaitGroup

	// Add 1 to the wait group
	wg.Add(1)

	// Start the agent in a goroutine
	go func() {
		defer wg.Done() // When the agent.Run() exits, signal that it's done
		a.Run(ctx)
	}()

	// Start the ConsumeOutMessages and ConsumeInMessages goroutines
	go a.ConsumeOutMessages()

	log.Printf("YOOOHOO")
	time.Sleep(1 * time.Second)

	// Send a message
	a.OutMessageQueue <- gAgents.Message{
		Receiver: "0.0.0.0:8002",
		Sender:   "0.0.0.0:8003",
		Type:     "greet",
		Content:  "Hello!",
	}
	log.Printf("Added the message in the Queue")

	// Wait for the wait group to signal that the agent.Run() is done
	wg.Wait()

	log.Printf("The len of the OutMessageQueue is : %v", len(a.OutMessageQueue))
}
