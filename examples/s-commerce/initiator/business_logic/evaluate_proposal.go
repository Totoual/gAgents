package businesslogic

import (
	"fmt"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/s-commerce/protocol"
)

type EvaluateProposoal struct{}

func (d *EvaluateProposoal) Apply(m gAgents.Message) gAgents.Message {
	message, ok := m.(*protocol.ProposalResponse)
	if !ok {
		fmt.Errorf("expected a CustomMessage")
	}

	var response_message gAgents.Message

	if message.Content.Price <= 35.0 {
		fmt.Println(message)

		header := protocol.MessageHeader{
			Protocol:       "Accept Negotiation",
			Performative:   protocol.ACCEPT,
			ConversationID: message.ConversationID,
			InReplyWith:    message.InReplyWith,
			Receiver:       message.Sender,
			Sender:         message.Receiver,
		}

		response_message = &protocol.AcceptanceMessage{
			MessageHeader: header,
		}

	} else {

		header := protocol.MessageHeader{
			Protocol:       "Reject Negotiation",
			Performative:   protocol.REJECT,
			ConversationID: message.ConversationID,
			InReplyWith:    message.InReplyWith,
			Receiver:       message.Sender,
			Sender:         message.Receiver,
		}

		response_message = &protocol.RejectionMessage{
			MessageHeader: header,
		}
	}

	return response_message
}
