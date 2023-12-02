package fipa

import (
	"fmt"

	gAgents "github.com/totoual/gAgents/agent"
)

type UniversalHandler struct {
	agent          *gAgents.Agent
	businessLogics map[Performative]BusinessLogic
}

func NewUniversalHandler() *UniversalHandler {
	return &UniversalHandler{
		businessLogics: make(map[Performative]BusinessLogic),
	}
}

func (h *UniversalHandler) SetBusinessLogic(performative Performative, logic BusinessLogic) {
	h.businessLogics[performative] = logic
}

func (h *UniversalHandler) HandleMessage(envelope gAgents.Envelope) {
	m, err := envelope.ToMessage(&FIPAMessage{})

	if err != nil {
		fmt.Println("Error deserializing message:", err)
		return
	}

	fipaMsg, ok := m.(*FIPAMessage)
	if !ok {
		fmt.Errorf("Couldn't type cast")
		return
	}

	switch m.GetPerformative() {
	case int(CFP):
		h.handleCFP(fipaMsg)
	case int(PROPOSAL):
		h.handleProposal(fipaMsg)
	case int(ACCEPT):
		handleAccept(fipaMsg)
	case int(REJECT):
		handleReject(fipaMsg)
	case int(TERMINATE):
		handleTerminate(fipaMsg)
	case int(COMPLETE):
		handleComplete(fipaMsg)
	case int(INFORM):
		handleInform(fipaMsg)
	default:
		fmt.Println("Unknown performative")
	}
}

func (h *UniversalHandler) handleCFP(fipaMsg *FIPAMessage) {
	logic, exists := h.businessLogics[CFP]
	if !exists {
		// Handle missing logic
		return
	}

	response := logic.Apply(fipaMsg.Performative, &fipaMsg.Content)
	var content FipaContent
	switch response.Performative {
	case CFP:
		content = response.AdditionalInfo.(CFPContentType) // Type assertion
	case PROPOSAL:
		content = response.AdditionalInfo.(ProposalCFPContentType)
		// ... other performatives
	}

	responseMsg := NewFIPAMessage(
		fipaMsg.Nonce+1,
		fipaMsg.Protocol,
		response.Performative,
		fipaMsg.ConversationID,
		fipaMsg.ReplyWith,
		fipaMsg.ReplyBy,
		fipaMsg.InReplyWith,
		fipaMsg.Sender,
		fipaMsg.Receiver,
		content)

	// Send the response message
	h.agent.OutMessageQueue <- *gAgents.NewEnvelope(responseMsg)
}

func (h *UniversalHandler) handleProposal(fipaMsg *FIPAMessage) {

}

func handleAccept(m *FIPAMessage) {

}

func handleReject(m *FIPAMessage) {

}

func handleTerminate(m *FIPAMessage) {

}

func handleComplete(m *FIPAMessage) {

}

func handleInform(m *FIPAMessage) {

}
