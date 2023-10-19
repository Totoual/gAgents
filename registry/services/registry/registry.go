package registry

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/timestamp"
	gAgents "github.com/totoual/gAgents/agent"
	pb "github.com/totoual/gAgents/protos/registry"
	"google.golang.org/grpc"
)

type RegistryService struct {
	pb.UnimplementedAgentRegistryServer
	Agents          map[string]*pb.AgentRegistration // Swap this to a real db.
	eventDispatcher *gAgents.EventDispatcher
}

func (rs *RegistryService) Init(srv *grpc.Server) {
	pb.RegisterAgentRegistryServer(srv, rs)
}

func NewRegistryService(ed *gAgents.EventDispatcher, ts *gAgents.TaskScheduler) *RegistryService {
	rs := &RegistryService{
		Agents:          make(map[string]*pb.AgentRegistration),
		eventDispatcher: ed,
	}

	return rs
}

func (rs *RegistryService) RegisterAgent(ctx context.Context, a *pb.AgentRegistration) (*pb.RegistrationResponse, error) {
	// For simplicity, using a map (agents) to store agents in-memory. In a real application, you'd likely want to store this data in a database.

	// Check if agent already exists
	if _, exists := rs.Agents[a.UniqueId]; exists {
		return &pb.RegistrationResponse{
			Success: false,
			Message: "Agent with this ID already exists.",
		}, nil
	}

	// Store agent details
	rs.Agents[a.UniqueId] = a

	// Log for debugging
	log.Printf("Registered agent: %v", a.UniqueId)

	return &pb.RegistrationResponse{
		Success: true,
		Message: "Agent registered successfully.",
		Topics:  []string{"negotiation", "services"},
	}, nil

}

func (rs *RegistryService) UnregisterAgent(ctx context.Context, a *pb.AgentUnregistration) (*pb.UnregistrationResponse, error) {
	// Check if agent exists
	if _, exists := rs.Agents[a.UniqueId]; !exists {
		return &pb.UnregistrationResponse{
			Success: false,
			Message: "Agent with this ID does not exist.",
		}, nil
	}

	// Remove agent from the map
	delete(rs.Agents, a.UniqueId)

	// Log for debugging
	log.Printf("Unregistered agent: %v", a.UniqueId)

	return &pb.UnregistrationResponse{
		Success: true,
		Message: "Agent unregistered successfully.",
	}, nil
}

func (rs *RegistryService) SendHeartbeat(ctx context.Context, hb *pb.Heartbeat) (*pb.HealthCheckResponse, error) {
	// Look up the agent by unique ID.
	agent, exists := rs.Agents[hb.UniqueId]
	if !exists {
		return &pb.HealthCheckResponse{
			Success: false,
			Message: "No agent found with the provided unique ID",
		}, nil
	}

	// Update the last heartbeat timestamp for the agent.
	agent.LastHeartbeat = &timestamp.Timestamp{Seconds: hb.Timestamp.Seconds, Nanos: hb.Timestamp.Nanos}

	// You could also update a database, if you are using one, to store the last heartbeat timestamp.

	log.Printf("Received heartbeat from agent: %v at %v", hb.UniqueId, hb.Timestamp.AsTime())

	return &pb.HealthCheckResponse{
		Success: true,
		Message: "Heartbeat received successfully",
	}, nil
}

func (rs *RegistryService) Search(ctx context.Context, sm *pb.SearchMessage) (*pb.SearchMessageResponse, error) {
	// For simplicity, just returning a success message without actually doing anything
	return &pb.SearchMessageResponse{
		Success: true,
		Message: "Search request received.",
	}, nil
}

/*
Create functions to enable search functionality. The search functionaity will
create an event for the kafka.sendMessage function so we can inform the rest
of the agents about the search.
*/
