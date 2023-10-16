package registry

import (
	"context"
	"log"

	gAgents "github.com/totoual/gAgents/agent"
	pb "github.com/totoual/gAgents/registry/proto"
	"google.golang.org/grpc"
)

const AgentRegisteredEventType = "AgentRegistered"

type RegistryService struct {
	pb.UnimplementedAgentRegistryServer
	agents          map[string]*pb.AgentRegistration // Swap this to a real db.
	eventDispatcher *gAgents.EventDispatcher
}

func (rs *RegistryService) Init(srv *grpc.Server) {
	pb.RegisterAgentRegistryServer(srv, rs)
}

func NewRegistryService(ed *gAgents.EventDispatcher) *RegistryService {
	return &RegistryService{
		agents:          make(map[string]*pb.AgentRegistration),
		eventDispatcher: ed,
	}

}

func (rs *RegistryService) RegisterAgent(ctx context.Context, a *pb.AgentRegistration) (*pb.RegistrationResponse, error) {
	// For simplicity, using a map (agents) to store agents in-memory. In a real application, you'd likely want to store this data in a database.

	// Check if agent already exists
	if _, exists := rs.agents[a.UniqueId]; exists {
		return &pb.RegistrationResponse{
			Success: false,
			Message: "Agent with this ID already exists.",
		}, nil
	}

	// Store agent details
	rs.agents[a.UniqueId] = a

	// Log for debugging
	log.Printf("Registered agent: %v", a.UniqueId)

	// Publish a successfully registration event.
	event := gAgents.Event{
		Type:    AgentRegisteredEventType,
		Payload: a,
	}
	rs.eventDispatcher.Publish(event)

	return &pb.RegistrationResponse{
		Success: true,
		Message: "Agent registered successfully.",
	}, nil
}
