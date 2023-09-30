package businesslogic

import (
	"fmt"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/s-commerce/protocol"
)

type DiscountLogic struct{}

func (d *DiscountLogic) Apply(m gAgents.Message) gAgents.Message {
	message, ok := m.(*protocol.CFPRequest)
	if !ok {
		fmt.Errorf("expected a CustomMessage")
	}

	total := message.Content.Order.CalculateTotalCost()

	var reply_message gAgents.Message

	if float64(message.Content.MaxPrice) < total {
		reply_message := protocol.RejectionMessage{
			MessageHeader: protocol.MessageHeader{
				Protocol:       "Reject Negotiation",
				Performative:   protocol.REJECT,
				ConversationID: message.ConversationID,
				InReplyWith:    message.ReplyWith,
				Receiver:       message.Sender,
				Sender:         message.Receiver,
			},
		}
		return reply_message
	} else {
		fmt.Printf("Total price is %f\n", total)
		proposal := protocol.ProposalContent{
			Order: message.Content.Order,
			Price: total - total*0.01,
		}

		reply_message := protocol.ProposalResponse{
			MessageHeader: protocol.MessageHeader{
				Protocol:       "Proposal Negotiation",
				Performative:   protocol.PROPOSAL,
				ConversationID: message.ConversationID,
				InReplyWith:    message.ReplyWith,
				Receiver:       message.Sender,
				Sender:         message.Receiver,
			},
			Content: proposal,
		}
	}

	return reply_message // Construct and return the appropriate response message
}
