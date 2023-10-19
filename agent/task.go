package gAgents

import (
	"context"
	"log"
	"time"
)

type Task interface {
	ID() string
	ScheduledAt() time.Time
	Execute()
	Interval() time.Duration // Return 0 if the task is not recurring
	StopCondition() bool     // Return true if the task should stop
	RescheduleTaskAt(time.Time) Task
}

type TaskScheduler struct {
	Tasks map[string]Task
	ctx   context.Context
}

func NewTaskScheduler(c context.Context) *TaskScheduler {
	return &TaskScheduler{
		Tasks: make(map[string]Task),
		ctx:   c,
	}
}

func (ts *TaskScheduler) AddTask(task Task) {
	ts.Tasks[task.ID()] = task
}

func (ts *TaskScheduler) RemoveTask(taskID string) {
	delete(ts.Tasks, taskID)
}

func (ts *TaskScheduler) ExecuteTasks() {
	ticker := time.NewTicker(time.Second * 10) // check every 10 seconds, adjust as needed
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentTime := time.Now().UTC()

			for taskID, task := range ts.Tasks {
				if task.ScheduledAt().Before(currentTime) || task.ScheduledAt().Equal(currentTime) {
					// Execute the task
					log.Printf("Executing the task!")
					task.Execute()

					if task.Interval() == 0 {
						ts.RemoveTask(taskID)
						continue
					}

					// Check if the task has a recurrence interval
					interval := task.Interval()
					if interval > 0 {
						// Schedule the next execution
						taskID := task.ID()
						nextExecution := currentTime.Add(interval)
						ts.Tasks[taskID] = task.RescheduleTaskAt(nextExecution)
					} else if task.StopCondition() {
						// If the task meets its stop condition, remove it from the scheduler
						ts.RemoveTask(taskID)
					}
				}
			}
		}
	}
}
