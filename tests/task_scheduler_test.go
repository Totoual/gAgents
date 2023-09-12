package tests_test

import (
	"context"
	"testing"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
)

type TestTask struct {
	id          string
	scheduledAt time.Time
	executed    bool
}

func (t *TestTask) ID() string {
	return t.id
}

func (t *TestTask) Type() string {
	return "TestTask"
}

func (t *TestTask) ScheduledAt() time.Time {
	return t.scheduledAt
}

func (t *TestTask) Parameters() map[string]interface{} {
	return nil
}

func (t *TestTask) Execute() {
	t.executed = true
}

func (t *TestTask) Interval() time.Duration {
	return 0
}

func (t *TestTask) StopCondition() bool {
	return t.executed
}

func (t *TestTask) RescheduleTaskAt(newTime time.Time) gAgents.Task {
	return &TestTask{
		id:          t.id,
		scheduledAt: newTime,
	}
}

func TestTaskScheduler_ExecuteTasks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	taskScheduler := gAgents.NewTaskScheduler(ctx)

	// Create a test task and add it to the scheduler
	now := time.Now()
	task := &TestTask{
		id:          "test-task",
		scheduledAt: now,
	}

	taskScheduler.AddTask(task)

	// Execute tasks
	taskScheduler.ExecuteTasks()

	// Check if the task was executed
	if !task.executed {
		t.Errorf("Task was not executed")
	}

	// Check if the task was removed from the scheduler
	if _, exists := taskScheduler.Tasks[task.ID()]; exists {
		t.Errorf("Task was not removed from the scheduler")
	}

	cancel()
}
