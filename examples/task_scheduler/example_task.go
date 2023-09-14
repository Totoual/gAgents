package main

import (
	"context"
	"fmt"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
)

type ExampleTask struct {
	id          string
	scheduledAt time.Time
	interval    time.Duration
}

func (t *ExampleTask) ID() string {
	return t.id
}

func (t *ExampleTask) Type() string {
	return "example"
}

func (t *ExampleTask) ScheduledAt() time.Time {
	return t.scheduledAt
}

func (t *ExampleTask) Parameters() map[string]interface{} {
	return nil // Implement if needed
}

func (t *ExampleTask) Execute() {
	fmt.Printf("Executing task %s\n", t.id)
}

func (t *ExampleTask) Interval() time.Duration {
	return t.interval
}

func (t *ExampleTask) StopCondition() bool {
	// Implement stop condition logic if needed
	return false
}

func (t *ExampleTask) RescheduleTaskAt(newTime time.Time) gAgents.Task {
	fmt.Printf("Rescheduling task %s\n", t.id)
	return &ExampleTask{
		id:          t.id,
		scheduledAt: newTime,
		interval:    t.interval,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a task scheduler
	taskScheduler := gAgents.NewTaskScheduler(ctx)

	// Create an example task
	exampleTask := &ExampleTask{
		id:          "task1",
		scheduledAt: time.Now().Add(5 * time.Second),
		interval:    10 * time.Second, // Set interval for recurring task
	}

	// Add the task to the scheduler
	taskScheduler.AddTask(exampleTask)

	// Execute tasks
	go func() {
		for {
			taskScheduler.ExecuteTasks()
			time.Sleep(time.Second) // Adjust the sleep time as needed
		}
	}()

	// Keep the program running for a while to see task execution
	time.Sleep(30 * time.Second)
}
