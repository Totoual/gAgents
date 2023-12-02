package fipa

import (
	"encoding/json"
)

type Performative int

type FipaContent interface{}

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

type FIPAMessage struct {
	Nonce          int          `json:"nonce"`
	Protocol       string       `json:"protocol"`
	Performative   Performative `json:"performative"`
	ConversationID string       `json:"conversation-id"`
	ReplyWith      string       `json:"reply-with"`
	ReplyBy        string       `json:"reply-by"`
	InReplyWith    string       `json:"in-reply-with,omitempty"`
	Receiver       string       `json:"receiver"`
	Sender         string       `json:"sender"`
	Content        FipaContent  `json:"content"`
}

func NewFIPAMessage(
	nonce int,
	protocol string,
	performative Performative,
	converstaion_id string,
	reply_with string,
	reply_by string,
	in_reply_with string,
	receiver string,
	sender string,
	content FipaContent) *FIPAMessage {
	return &FIPAMessage{
		Nonce:          nonce,
		Protocol:       protocol,
		Performative:   performative,
		ConversationID: converstaion_id,
		ReplyWith:      reply_with,
		ReplyBy:        reply_by,
		InReplyWith:    in_reply_with,
		Receiver:       receiver,
		Sender:         sender,
		Content:        content,
	}
}

func (fm *FIPAMessage) GetProtocol() string {
	return fm.Protocol
}

func (fm *FIPAMessage) GetReceiver() string {
	return fm.Receiver
}

func (fm *FIPAMessage) GetSender() string {
	return fm.Sender
}

func (fm *FIPAMessage) Serialize() ([]byte, error) {
	return json.Marshal(fm)
}

func (fm *FIPAMessage) GetPerformative() int {
	return int(fm.Performative)
}
