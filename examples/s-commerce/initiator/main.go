package main

import (
	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/s-commerce/initiator/acts"
	businesslogic "github.com/totoual/gAgents/examples/s-commerce/initiator/business_logic"
	"github.com/totoual/gAgents/examples/s-commerce/protocol"
)

func main() {
	agent := gAgents.NewAgent("Alice", "127.0.0.1:8000")
	cfpAct := &acts.SendMessageAct{}
	agent.RegisterAct(cfpAct)

	l := &businesslogic.EvaluateProposoal{}

	handler := protocol.NewProposalHandler(agent, l)

	agent.RegisterHandler(handler.Name, handler)
	agent.Run()

}
