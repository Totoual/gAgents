package registry

import (
	"context"
	"log"

	pb "github.com/totoual/gAgents/registry/proto"
	"google.golang.org/grpc"
)

type RegistryService struct {
	pb.UnimplementedAgentRegistryServer
	agents map[string]*pb.AgentRegistration // Swap this to a real db.
}

func (rs *RegistryService) Init(srv *grpc.Server) {
	pb.RegisterAgentRegistryServer(srv, rs)
}

func NewRegistryService() *RegistryService {
	return &RegistryService{
		agents: make(map[string]*pb.AgentRegistration),
	}

}

func (rs *RegistryService) RegisterAgent(ctx context.Context, agent *pb.AgentRegistration) (*pb.RegistrationResponse, error) {
	// For simplicity, using a map (agents) to store agents in-memory. In a real application, you'd likely want to store this data in a database.

	// Check if agent already exists
	if _, exists := rs.agents[agent.UniqueId]; exists {
		return &pb.RegistrationResponse{
			Success: false,
			Message: "Agent with this ID already exists.",
		}, nil
	}

	// Store agent details
	rs.agents[agent.UniqueId] = agent

	// Log for debugging
	log.Printf("Registered agent: %v", agent.UniqueId)

	return &pb.RegistrationResponse{
		Success: true,
		Message: "Agent registered successfully.",
	}, nil
}
