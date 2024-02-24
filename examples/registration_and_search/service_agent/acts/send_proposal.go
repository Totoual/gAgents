package acts

import (
	"encoding/json"
	"log"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/registration_and_search/service_agent/config"
	"github.com/totoual/gAgents/protocols/fipa"
	"github.com/totoual/gAgents/services/kafka"
)

type SendProposalAct struct {
	Config config.AgentConfig
	ed     *gAgents.EventDispatcher
	agent  *gAgents.Agent
}

func NewSendProposalAct(config config.AgentConfig, ed *gAgents.EventDispatcher, a *gAgents.Agent) *SendProposalAct {
	act := &SendProposalAct{
		Config: config,
		ed:     ed,
		agent:  a,
	}
	ed.Subscribe(kafka.KafkaSearch, act.handleKafkaSearchEvent)
	return act
}

func (sp *SendProposalAct) handleKafkaSearchEvent(e gAgents.Event) {
	log.Println("Received an event in search event! Building the proposal")
	log.Println(e.Type)
	var msg kafka.KafkaConsumerMessage
	if err := json.Unmarshal(e.Payload.([]byte), &msg); err != nil {
		log.Printf("Could not unmarshal the message: %v\n", err)
		return
	}

	log.Println(msg)
	proposalMsg := fipa.NewFIPAMessage(
		0,
		"fipa",
		fipa.PROPOSAL,
		"conversation_id_value",
		"reply_with_value",
		msg.GrpcAddress,
		sp.agent.Addr,
		fipa.FipaProposal{})
	log.Println("Proposal message: ", proposalMsg)

	sp.agent.OutMessageQueue <- *gAgents.NewEnvelope(proposalMsg)

}
