package agent

import "fmt"

type Message struct {
	Sender  string
	Content string
}

type Agent struct {
	name    string
	inQueue chan Message
	exit    chan struct{}
}

func NewAgent(name string) *Agent {
	return &Agent{
		name:    name,
		inQueue: make(chan Message),
		exit:    make(chan struct{}),
	}
}

func (a *Agent) SendMessage(message Message, recipient *Agent) {
	recipient.ReceiveMessage(message)
}

func (a *Agent) ReceiveMessage(message Message) {
	a.inQueue <- message
}

func (a *Agent) Run() {
	for {
		select {
		case message := <-a.inQueue:
			// Process the incoming message
			fmt.Printf("[%s] Received message from %s: %s\n", a.name, message.Sender, message.Content)
		case <-a.exit:
			// Received exit signal, clean up and exit the loop
			close(a.inQueue)
			return
		}
	}
}

func (a *Agent) Exit() {
	// Send an exit signal to the agent's exit channel
	close(a.exit)
}
