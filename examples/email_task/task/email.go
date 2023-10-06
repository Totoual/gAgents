package task

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	gAgents "github.com/totoual/gAgents/agent"
	"gopkg.in/gomail.v2"
)

type EmailTask struct {
	id          string
	scheduledAt time.Time
	interval    time.Duration
	sender      string
	receiver    string
	subject     string
	body        string
	password    string
	stop        bool
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

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	password := os.Getenv("gmail_code")
	log.Printf(password)
	return &EmailTask{
		id:          uuid.New().String(),
		scheduledAt: sendTime,
		interval:    time.Duration(interval),
		sender:      sender,
		receiver:    receiver,
		subject:     subject,
		body:        body,
		password:    password,
		stop:        false,
	}

}

func (t *EmailTask) ID() string {
	return t.id
}

func (t *EmailTask) ScheduledAt() time.Time {
	return t.scheduledAt
}

func (t *EmailTask) Execute() {
	fmt.Printf("Executing task %s\n", t.id)
	message := gomail.NewMessage()
	message.SetHeader("From", t.sender)
	message.SetHeader("To", t.receiver)
	message.SetHeader("Subject", t.subject)
	message.SetBody("text/plain", t.body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, t.sender, t.password)

	if err := dialer.DialAndSend(message); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Email sent!")
	}

	t.stop = true
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
