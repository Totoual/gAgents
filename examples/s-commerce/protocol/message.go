package protocol

import (
	"encoding/json"

	gAgents "github.com/totoual/gAgents/agent"
)

type Performative int

type BusinessLogic interface {
	Apply(gAgents.Message) gAgents.Message
}

const (
	CFP Performative = iota
	PROPOSAL
	ACCEPT
	REJECT
	TERMINATE
	COMPLETE
	INFORM

	// Add more performative constants here
)

type MessageHeader struct {
	Protocol       string       `json:"protocol"`
	Performative   Performative `json:"performative"`
	ConversationID string       `json:"conversation-id"`
	ReplyWith      string       `json:"reply-with"`
	ReplyBy        string       `json:"reply-by"`
	InReplyWith    string       `json:"in-reply-with,omitempty"`
	Receiver       string       `json:"receiver"`
	Sender         string       `json:"sender"`
}

func (m *MessageHeader) GetProtocol() string {
	return m.Protocol
}

func (m *MessageHeader) GetReceiver() string {
	return m.Receiver
}

func (m *MessageHeader) GetSender() string {
	return m.Sender
}

func (m *MessageHeader) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MessageHeader) GetPerformative() int {
	return int(m.Performative)
}

type Item struct {
	SKU    string  `json:"sku"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
}

type Order struct {
	Items []Item `json:"items"`
}

func (o Order) CalculateTotalCost() float64 {
	var total float64
	for _, item := range o.Items {
		total += (item.Price * float64(item.Amount))
	}
	return total
}

type CFPRequest struct {
	MessageHeader
	Content CFPContent `json:"content"`
}

type CFPContent struct {
	Order    Order `json:"order"`
	MaxPrice int   `json:"max-price"`
}

type ProposalResponse struct {
	MessageHeader
	Content ProposalContent `json:"content"`
}

type ProposalContent struct {
	Order Order   `json:"order"`
	Price float64 `json:"price"`
}

type AcceptanceMessage struct {
	MessageHeader
	Content Content `json:"content"`
}

type RejectionMessage struct {
	MessageHeader
	Content Content `json:"content"`
}

type Content struct{}

type TerminationMessage struct {
	MessageHeader
	Content TerminationContent `json:"content"`
}

type TerminationContent struct {
	Reason string `json:"reason"`
	Note   string `json:"note"`
}

type CompletionMessage struct {
	MessageHeader
	Content CompletionContent `json:"content"`
}

type CompletionContent struct {
	Reason string `json:"reason"`
	// Other relevant fields specific to the completion
}
