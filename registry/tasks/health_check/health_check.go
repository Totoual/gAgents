package healthcheck

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	gAgents "github.com/totoual/gAgents/agent"
	pb "github.com/totoual/gAgents/protos/registry"
	"github.com/totoual/gAgents/registry/services/registry"
)

type HeartbeatTask struct {
	id              string
	scheduledAt     time.Time
	interval        time.Duration
	registryService *registry.RegistryService
}

func NewHeartbeatTask(st time.Time, it time.Duration, rs *registry.RegistryService) *HeartbeatTask {
	return &HeartbeatTask{
		id:              uuid.New().String(),
		scheduledAt:     st,
		interval:        it,
		registryService: rs,
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
	fmt.Printf("Executing task %s\n", t.id)
	for id, agent := range t.registryService.Agents {
		if agent == nil || agent.LastHeartbeat == nil {
			log.Printf("Warning: agent or LastHeartbeat is nil for agent ID %v", id)
			continue // Skip to next iteration
		}
		// Assuming a 5-minute timeout for simplicity
		if time.Since(agent.LastHeartbeat.AsTime()) > 5*time.Minute {
			// Unregister agent if heartbeat is too old
			_, err := t.registryService.UnregisterAgent(context.Background(), &pb.AgentUnregistration{UniqueId: id})
			if err != nil {
				log.Printf("Error unregistering agent %v: %v", id, err)
			}
		}

	}
}

func (t *HeartbeatTask) Interval() time.Duration {
	return t.interval
}

func (t *HeartbeatTask) StopCondition() bool {
	// Implement stop condition logic if needed
	return false
}

func (t *HeartbeatTask) RescheduleTaskAt(newTime time.Time) gAgents.Task {
	fmt.Printf("Rescheduling task %s\n", t.id)
	return &HeartbeatTask{
		id:              t.id,
		scheduledAt:     newTime,
		interval:        t.interval,
		registryService: t.registryService,
	}
}
