package main

import (
	"context"

	"github.com/totoual/gAgents"
)

func main() {

	a := gAgents.NewAgent("Test", "0.0.0.0:8000")
	ctx := context.Background()
	go a.ConsumeInMessages()
	a.Run(ctx)

}
