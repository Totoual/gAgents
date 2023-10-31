package tasks

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/registration_and_search/service_agent/config"
	pb "github.com/totoual/gAgents/protos/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HeartbeatTask struct {
	id               string
	scheduledAt      time.Time
	interval         time.Duration
	Config           config.AgentConfig
	stop_contidition bool
}

func NewHeartbeatTask(st time.Time, it time.Duration, config config.AgentConfig) *HeartbeatTask {
	return &HeartbeatTask{
		id:               uuid.New().String(),
		scheduledAt:      st,
		interval:         it,
		Config:           config,
		stop_contidition: false,
	}
}

func (t *HeartbeatTask) ID() string {
	return t.id
}

func (t *HeartbeatTask) Type() string {
	return "Repeated"
}

func (t *HeartbeatTask) ScheduledAt() time.Time {
	return t.scheduledAt
}

func (t *HeartbeatTask) Execute() {
	conn, err := grpc.Dial(t.Config.RegistryURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewAgentRegistryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// Send the message via gRPC
	response, err := client.SendHeartbeat(ctx, &pb.Heartbeat{
		UniqueId:  t.Config.UniqueID,
		Timestamp: timestamppb.New(time.Now()),
	})
	if err != nil {
		log.Println(err)
	}
	if !response.Success {
		t.stop_contidition = true
	}
	log.Println(response)
}

func (t *HeartbeatTask) Interval() time.Duration {
	return t.interval
}

func (t *HeartbeatTask) StopCondition() bool {
	// Implement stop condition logic if needed
	return t.stop_contidition
}

func (t *HeartbeatTask) RescheduleTaskAt(newTime time.Time) gAgents.Task {
	fmt.Printf("Rescheduling task %s\n", t.id)
	return &HeartbeatTask{
		id:          t.id,
		scheduledAt: newTime,
		interval:    t.interval,
		Config:      t.Config,
	}
}
