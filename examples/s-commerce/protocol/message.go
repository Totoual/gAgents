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

// CFP (Call for Proposal)

type CFPRequest struct {
	Protocol       string       `json:"protocol"`
	Performative   Performative `json:"performative"`
	ConversationID string       `json:"conversation-id"`
	ReplyWith      string       `json:"reply-with"`
	Content        CFPContent   `json:"content"`
	ReplyBy        string       `json:"reply-by"`
	Receiver       string       `json:"receiver"`
	Sender         string       `json:"sender"`
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

type CFPContent struct {
	Order    Order `json:"order"`
	MaxPrice int   `json:"max-price"`
}

func (cfp CFPRequest) GetProtocol() string {
	return cfp.Protocol
}

func (cfp CFPRequest) GetReceiver() string {
	return cfp.Receiver
}

func (cfp CFPRequest) GetSender() string {
	return cfp.Sender
}

func (cfp CFPRequest) Serialize() ([]byte, error) {
	return json.Marshal(cfp)
}

func (cfp CFPRequest) GetPerformative() int {
	return int(cfp.Performative)
}

// Receiving Proposals

type ProposalResponse struct {
	Protocol       string          `json:"protocol"`
	Performative   Performative    `json:"performative"`
	ConversationID string          `json:"conversation-id"`
	ReplyWith      string          `json:"reply-with"`
	InReplyWith    string          `json:"in-reply-with"`
	Content        ProposalContent `json:"content"`
	ReplyBy        string          `json:"reply-by"`
	Receiver       string          `json:"receiver"`
	Sender         string          `json:"sender"`
}

type ProposalContent struct {
	Order Order   `json:"order"`
	Price float64 `json:"price"`
}

func (p ProposalResponse) GetProtocol() string {
	return p.Protocol
}

func (p ProposalResponse) GetReceiver() string {
	return p.Receiver
}

func (p ProposalResponse) GetSender() string {
	return p.Sender
}

func (p ProposalResponse) Serialize() ([]byte, error) {
	return json.Marshal(p)
}

func (p ProposalResponse) GetPerformative() int {
	return int(p.Performative)
}

// Proposal Evaluation
type Content struct{}

// Sending Acceptance or Rejection Messages
type AcceptanceMessage struct {
	Protocol       string       `json:"protocol"`
	Performative   Performative `json:"performative"`
	ConversationID string       `json:"conversation-id"`
	InReplyWith    string       `json:"in-reply-with"`
	ReplyBy        string       `json:"reply-by"`
	Receiver       string       `json:"receiver"`
	Sender         string       `json:"sender"`
	Content        Content      `json:"content"`
}

func (a AcceptanceMessage) GetProtocol() string {
	return a.Protocol
}

func (a AcceptanceMessage) GetReceiver() string {
	return a.Receiver
}

func (a AcceptanceMessage) GetSender() string {
	return a.Sender
}

func (a AcceptanceMessage) Serialize() ([]byte, error) {
	return json.Marshal(a)
}

func (a AcceptanceMessage) GetPerformative() int {
	return int(a.Performative)
}

type RejectionMessage struct {
	Protocol       string       `json:"protocol"`
	Performative   Performative `json:"performative"`
	ConversationID string       `json:"conversation-id"`
	InReplyWith    string       `json:"in-reply-with"`
	ReplyBy        string       `json:"reply-by"`
	Receiver       string       `json:"receiver"`
	Sender         string       `json:"sender"`
	Content        Content      `json:"content"`
}

func (r RejectionMessage) GetProtocol() string {
	return r.Protocol
}

func (r RejectionMessage) GetReceiver() string {
	return r.Receiver
}

func (r RejectionMessage) GetSender() string {
	return r.Sender
}

func (r RejectionMessage) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

func (r RejectionMessage) GetPerformative() int {
	return int(r.Performative)
}

// Completion or Termination
type TerminationMessage struct {
	Protocol       string             `json:"protocol"`
	Performative   Performative       `json:"performative"`
	ConversationID string             `json:"conversation-id"`
	Content        TerminationContent `json:"content"`
	Receiver       string             `json:"receiver"`
	Sender         string             `json:"sender"`
}

type TerminationContent struct {
	Reason string `json:"reason"`
	Note   string `json:"note"`
}

func (t TerminationMessage) GetProtocol() string {
	return t.Protocol
}

func (t TerminationMessage) GetReceiver() string {
	return t.Receiver
}

func (t TerminationMessage) GetSender() string {
	return t.Sender
}

func (t TerminationMessage) Serialize() ([]byte, error) {
	return json.Marshal(t)
}

func (t TerminationMessage) GetPerformative() int {
	return int(t.Performative)
}

type CompletionMessage struct {
	Protocol       string            `json:"protocol"`
	Performative   Performative      `json:"performative"`
	ConversationID string            `json:"conversation-id"`
	Content        CompletionContent `json:"content"`
	Receiver       string            `json:"receiver"`
	Sender         string            `json:"sender"`
}

type CompletionContent struct {
	Reason string `json:"reason"`
	// Other relevant fields specific to the completion
}

func (c CompletionMessage) GetProtocol() string {
	return c.Protocol
}

func (c CompletionMessage) GetReceiver() string {
	return c.Receiver
}

func (c CompletionMessage) GetSender() string {
	return c.Sender
}

func (c CompletionMessage) GetPerformative() int {
	return int(c.Performative)
}

func (c CompletionMessage) Serialize() ([]byte, error) {
	return json.Marshal(c)
}
