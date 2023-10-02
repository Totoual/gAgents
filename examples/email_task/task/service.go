package task

import (
	"context"

	gAgents "github.com/totoual/gAgents/agent"
	pb "github.com/totoual/gAgents/examples/email_task/generated"
	"google.golang.org/grpc"
)

type EmailService struct {
	pb.UnimplementedEmailServiceServer
	Scheduler *gAgents.TaskScheduler
}

func (es *EmailService) Init(srv *grpc.Server) {
	pb.RegisterEmailServiceServer(srv, es)
}

func NewEmailService(scheduler *gAgents.TaskScheduler) *EmailService {
	return &EmailService{
		Scheduler: scheduler,
	}
}

func (es *EmailService) CreateEmailTask(ctx context.Context, req *pb.EmailTask) (*pb.EmailTaskResponse, error) {
	// Use the provided data to create a new EmailTask.

	emailTask := NewEmailTask(
		req.ScheduleAt,
		int(req.Interval),
		req.Sender,
		req.Receiver,
		req.Subject,
		req.Body,
	)

	// Add the task to the TaskScheduler.
	es.Scheduler.AddTask(emailTask)

	return &pb.EmailTaskResponse{
		TaskId:  emailTask.ID(),
		Message: "Email task successfully created.",
	}, nil
}
