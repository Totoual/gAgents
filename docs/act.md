## Acts

Acts are a powerful feature of the gAgents framework, allowing you to define actions that an agent can perform at specific intervals or in response to certain events.

### How it Works

Acts are defined by implementing the `Act` interface, which includes methods for performing the action (`Perform()`) and specifying the interval at which the act should be executed (`GetInterval()`).

### Registering Acts

To make an agent perform an act, you simply register the act with the agent using `agent.RegisterAct(myAct)`.

### Performing Acts

Acts can be set to execute at regular intervals by specifying a non-zero interval in the `GetInterval()` method. If an interval is set, the act will be performed repeatedly at the specified interval.

### Stopping Acts

Acts can be stopped by either not specifying an interval (making it a one-time action) or by including a stop condition in the `Perform()` method.


### Example Usage

```go
// Define a custom act
package main

import (
	"fmt"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
)

// CustomAct is an example of a custom act.
type CustomAct struct {
	count int
}

// Perform is the method called when the act is performed.
func (a *CustomAct) Perform(agent *gAgents.Agent) {
	// Perform action...
	a.count++
	fmt.Printf("CustomAct performed. Count: %d\n", a.count)
}

// GetInterval returns the interval at which the act should be performed.
func (a *CustomAct) GetInterval() time.Duration {
	return time.Second * 2 // Perform every 2 seconds
}

func main() {
	// Create a new agent
	agent := gAgents.NewAgent("Alice", "localhost:8000")

	// Create a custom act
	myAct := &CustomAct{}

	// Register the act with the agent
	agent.RegisterAct(myAct)

	// Start performing acts
	go agent.PerformActs()

	// Let the agent run for a while to observe act execution
	time.Sleep(time.Second * 10)

	// Stop the agent (Optional)
	agent.Cancel()
}

```

In this example, we define a custom act CustomAct that increments a counter each time it is performed. The act is registered with the agent and set to execute every 2 seconds.

When you run this program, you should see output indicating that the act is being performed at regular intervals.