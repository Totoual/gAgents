package acts

import (
	"encoding/json"
	"log"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/registration_and_search/service_agent/config"
	"github.com/totoual/gAgents/services/kafka"
)

type SendProposalAct struct {
	Config config.AgentConfig
	ed     *gAgents.EventDispatcher
}

func NewSendProposalAct(config config.AgentConfig, ed *gAgents.EventDispatcher) *SendProposalAct {

	ed.Subscribe(kafka.KafkaSearch, handleKafkaSearchEvent)
	return &SendProposalAct{
		Config: config,
		ed:     ed,
	}
}

func handleKafkaSearchEvent(e gAgents.Event) {
	log.Println("Received an event")
	log.Println(e.Type)
	var msg kafka.KafkaConsumerMessage
	if err := json.Unmarshal(e.Payload.([]byte), &msg); err != nil {
		log.Printf("Could not unmarshal the message: %v\n", err)
		return
	}

	log.Println(msg)
}
