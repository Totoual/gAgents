package protocol

import (
	"fmt"
	"log"

	gAgents "github.com/totoual/gAgents/agent"
)

type CFPHandler struct {
	Name  string
	agent *gAgents.Agent
	logic BusinessLogic
}

func NewCFPHandler(a *gAgents.Agent, l BusinessLogic) *CFPHandler {
	return &CFPHandler{
		Name:  "CFP Negotiation",
		agent: a,
		logic: l,
	}
}

func (h *CFPHandler) HandleMessage(envelope gAgents.Envelope) {
	m, err := envelope.ToMessage(&CFPRequest{})
	if err != nil {
		fmt.Println("Error deserializing message:", err)
	}
	// Here needs to be the agent logic
	response_message := h.logic.Apply(m)
	if ok := response_message.GetPerformative() == int(PROPOSAL) || response_message.GetPerformative() == int(REJECT); !ok {
		log.Fatalf("The response message should be either of Performative Proposal or Reject")
	}

	// The new message needs to be of type Propose or Reject
	h.agent.OutMessageQueue <- *gAgents.NewEnvelope(response_message)

}

type ProposalHandler struct {
	Name  string
	agent *gAgents.Agent
	logic BusinessLogic
}

func NewProposalHandler(a *gAgents.Agent, l BusinessLogic) *ProposalHandler {
	return &ProposalHandler{
		Name:  "Proposal Negotiation",
		agent: a,
		logic: l,
	}
}

func (h *ProposalHandler) HandleMessage(envelope gAgents.Envelope) {
	m, err := envelope.ToMessage(&ProposalResponse{})
	if err != nil {
		fmt.Println("Error deserializing message:", err)
	}
	// Here needs to be the agent logic
	response_message := h.logic.Apply(m)
	if ok := response_message.GetPerformative() == int(PROPOSAL) ||
		response_message.GetPerformative() == int(REJECT) ||
		response_message.GetPerformative() == int(ACCEPT); !ok {

		log.Fatalf("The response message should be either of Performative Proposal or Reject")
	}

	// The new message needs to be of type Propose or Reject
	h.agent.OutMessageQueue <- *gAgents.NewEnvelope(response_message)

}

type AcceptanceHandler struct {
	Name  string
	agent *gAgents.Agent
	logic BusinessLogic
}

func NewAcceptanceHandler(a *gAgents.Agent) *AcceptanceHandler {
	return &AcceptanceHandler{
		Name:  "Accept Negotiation",
		agent: a,
	}
}

func (h *AcceptanceHandler) HandleMessage(envelope gAgents.Envelope) {
	m, err := envelope.ToMessage(&AcceptanceMessage{})
	if err != nil {
		fmt.Println("Error deserializing message:", err)
	}
	message, ok := m.(*AcceptanceMessage)
	if !ok {
		fmt.Errorf("expected a CustomMessage")
	}

	fmt.Println(message)

	// Here needs to be the agent logic
	// response_message := h.logic.Apply(m)
	// if ok := response_message.GetPerformative() == int(INFORM) ||
	// 	response_message.GetPerformative() == int(TERMINATE) ||
	// 	response_message.GetPerformative() == int(COMPLETE); !ok {

	// 	log.Fatalf("The response message should be either of Performative Proposal or Reject")
	// }

	// // The new message needs to be of type Propose or Reject
	// h.agent.OutMessageQueue <- *gAgents.NewEnvelope(response_message)

}
