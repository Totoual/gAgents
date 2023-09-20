package main

import (
	gAgents "github.com/totoual/gAgents/agent"
	businesslogic "github.com/totoual/gAgents/examples/s-commerce/participant/business_logic"
	"github.com/totoual/gAgents/examples/s-commerce/protocol"
)

func main() {
	a := gAgents.NewAgent("Bob", "127.0.0.1:8001")
	l := &businesslogic.DiscountLogic{}

	cfpHandler := protocol.NewCFPHandler(a, l)

	a.RegisterHandler(cfpHandler.Name, cfpHandler)
	acceptHandler := protocol.NewAcceptanceHandler(a)
	a.RegisterHandler(acceptHandler.Name, acceptHandler)
	a.Run()

}
