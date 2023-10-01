package acts

import (
	"time"

	"github.com/google/uuid"
	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/s-commerce/protocol"
)

type SendMessageAct struct{}

const (
	ProtocolName    = "CFP Negotiation"
	DefaultReceiver = "127.0.0.1:8001"
)

// Helper function to generate a new UUID string
func newUUID() string {
	return uuid.NewString()
}

// Helper function to create a new item
func newItem(amount int, price float64) protocol.Item {
	return protocol.Item{
		SKU:    newUUID(),
		Amount: amount,
		Price:  price,
	}
}

func (a *SendMessageAct) Perform(agent *gAgents.Agent) {
	header := protocol.MessageHeader{
		Protocol:       ProtocolName,
		Performative:   protocol.CFP,
		ConversationID: newUUID(),
		ReplyWith:      newUUID(),
		Receiver:       DefaultReceiver,
		Sender:         agent.Addr,
	}

	content := protocol.CFPContent{
		Order: protocol.Order{
			Items: []protocol.Item{
				newItem(5, 2.5),
			},
		},
		MaxPrice: 35,
	}

	message := &protocol.CFPRequest{
		MessageHeader: header,
		Content:       content,
	}

	agent.OutMessageQueue <- *gAgents.NewEnvelope(message)
}

func (a *SendMessageAct) GetInterval() time.Duration {
	return 0
}
