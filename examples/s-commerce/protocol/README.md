# Negotiation Protocol

This Go package provides a flexible framework for conducting negotiations between agents. The negotiation process includes a Call for Proposal (CFP), proposal responses, acceptance or rejection messages, and completion or termination messages.

## Key Components

### Performative

`Performative` is an enumeration that defines the different types of negotiation messages. The supported performative types are: CFP, PROPOSAL, ACCEPT, REJECT, TERMINATE, COMPLETE, and INFORM.

### Business Logic Interface

The `BusinessLogic` interface allows developers to define their own negotiation logic. The `Apply` method takes a message as input and returns a new message.

### Handlers

Handlers are components that process incoming messages based on their performative type. There are different handlers for CFPs, proposal responses, and acceptance messages.

- **CFPHandler**: Handles Call for Proposal (CFP) messages. It applies the defined negotiation logic and generates proposal or rejection responses.

- **ProposalHandler**: Handles proposal response messages. It applies the defined negotiation logic and generates new proposal, rejection, or acceptance responses.

- **AcceptanceHandler**: Handles acceptance messages. It currently prints the received message and provides a placeholder for implementing custom logic.

## Usage

1. **Define Custom Business Logic**: Implement the `BusinessLogic` interface to define your custom negotiation logic.

2. **Create Agents**: Create agents and specify their behavior, including message handlers and business logic.

3. **Register Handlers**: Register message handlers for different performative types to handle incoming messages appropriately.

4. **Conduct Negotiations**: Use the provided structs and methods to conduct negotiations between agents.

## Example

```go
package main

import (
    "github.com/totoual/gAgents/examples/s-commerce/protocol"
    gAgents "github.com/totoual/gAgents/agent"
)

func main() {
    a := gAgents.NewAgent("Bob", "127.0.0.1:8001")
    l := &DiscountLogic{}

    cfpHandler := protocol.NewCFPHandler(a, l)
    proposalHandler := protocol.NewProposalHandler(a, l)
    acceptanceHandler := protocol.NewAcceptanceHandler(a)

    a.RegisterHandler("CFP Negotiation", cfpHandler)
    a.RegisterHandler("Proposal Negotiation", proposalHandler)
    a.RegisterHandler("Accept Negotiation", acceptanceHandler)

    a.Run()
}
```