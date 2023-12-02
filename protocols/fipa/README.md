---

# FIPA Protocol Integration in gAgents

The FIPA protocol in gAgents allows for flexible and robust communication between agents using different performatives such as CFP (Call for Proposals), Proposals, Acceptance, and more. This guide will walk you through setting up the FIPA protocol with custom business logic.

## Files

- `fipa_message.go`: Defines the `FIPAMessage` struct and `FipaContent` interface used for all FIPA messages.
- `fipa_handler.go`: Contains the `UniversalHandler` for handling all FIPA messages.
- `fipa_business_logic_ctx.go`: Defines the `BusinessLogic` interface and `BusinessLogicContext` for implementing custom business logic.

## Setting Up

To use the FIPA protocol in your agent, you need to perform the following steps:

### 1. Define Your Content Types

Implement your own content types that fulfill the `FipaContent` interface. For example:

```go
type MyCustomContent struct {
    // Your fields here
}

func (mcc MyCustomContent) GetCart() *[]Item {
    // Implementation here
}
```

### 2. Implement Business Logic

Implement the `BusinessLogic` interface. Your implementation should handle the business logic specific to your application and return a `BusinessLogicContext`:

```go
type MyBusinessLogic struct {
    // Your fields and methods
}

func (mbl *MyBusinessLogic) Apply(performative Performative, content *FipaContent) *BusinessLogicContext {
    // Your logic here
    return &BusinessLogicContext{
        Performative:   // Your Performative,
        AdditionalInfo: // Your FipaContent or other data,
    }
}
```

### 3. Set Up the Universal Handler

In your agent's setup, initialize the `UniversalHandler` and set the business logic for each FIPA performative:

```go
agent := gAgents.NewAgent("AgentName", "localhost:8080")
handler := fipa.NewUniversalHandler(agent)
handler.SetBusinessLogic(fipa.CFP, &MyBusinessLogic{})
// Set other performatives as needed
```

### 4. Start Your Agent

Start your agent to begin handling messages:

```go
agent.Run()
```

## Message Flow

1. When a FIPA message is received, the `UniversalHandler` de-serializes it into a `FIPAMessage`.
2. The handler checks the performative of the message and calls the appropriate business logic.
3. The business logic processes the message and returns a `BusinessLogicContext`.
4. The handler constructs a new FIPA message based on the response and sends it out.

## Customization

- Implement different `BusinessLogic` for different scenarios or performatives.
- Define your own content types as per your application's requirements.
- Customize the `UniversalHandler` if additional message processing is needed.

---