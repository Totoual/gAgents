package task

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	gAgents "github.com/totoual/gAgents/agent"
)

type EmailTask struct {
	id          string
	scheduledAt time.Time
	interval    time.Duration
	sender      string
	receiver    string
	subject     string
	body        string
}

func NewEmailTask(
	scheduled_time string,
	interval int,
	sender string,
	receiver string,
	subject string,
	body string,
) *EmailTask {
	sendTime, err := time.Parse(time.RFC3339, scheduled_time)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return &EmailTask{
		id:          uuid.New().String(),
		scheduledAt: sendTime,
		interval:    time.Duration(interval),
		sender:      sender,
		receiver:    receiver,
		subject:     subject,
		body:        body,
	}

}

func (t *EmailTask) ID() string {
	return t.id
}

func (t *EmailTask) Type() string {
	return "example"
}

func (t *EmailTask) ScheduledAt() time.Time {
	return t.scheduledAt
}

func (t *EmailTask) Execute() {
	fmt.Printf("Executing task %s\n", t.id)

	// Send the email to the receiver.
}

func (t *EmailTask) Interval() time.Duration {
	return t.interval
}

func (t *EmailTask) StopCondition() bool {
	// Implement stop condition logic if needed
	return false
}

func (t *EmailTask) RescheduleTaskAt(newTime time.Time) gAgents.Task {
	fmt.Printf("Rescheduling task %s\n", t.id)
	return &EmailTask{
		id:          t.id,
		scheduledAt: newTime,
		interval:    t.interval,
		sender:      t.sender,
		receiver:    t.receiver,
		subject:     t.subject,
		body:        t.body,
	}
}
