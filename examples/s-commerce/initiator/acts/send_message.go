package acts

import (
	"time"

	"github.com/google/uuid"
	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/s-commerce/protocol"
)

type SendMessageAct struct{}

func (a *SendMessageAct) Perform(agent *gAgents.Agent) {
	message := protocol.CFPRequest{

		Protocol:       "CFP Negotiation",
		Performative:   protocol.CFP,
		ConversationID: uuid.NewString(),
		ReplyWith:      uuid.NewString(),
		Content: protocol.CFPContent{
			Order: protocol.Order{
				Items: []protocol.Item{
					{
						SKU:    uuid.NewString(),
						Amount: 5,
						Price:  2.5,
					},
				},
			},
			MaxPrice: 35,
		},
		Receiver: "127.0.0.1:8001",
		Sender:   agent.Addr,
	}

	agent.OutMessageQueue <- *gAgents.NewEnvelope(message)
}

func (a *SendMessageAct) GetInterval() time.Duration {
	return 0
}
