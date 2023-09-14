package main

import (
	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/server_client_communication/protocol"
)

func main() {

	a := gAgents.NewAgent("Test", "0.0.0.0:8002")
	a.RegisterHandler("greet", &protocol.GreetingHandler{})
	a.Run()

}
