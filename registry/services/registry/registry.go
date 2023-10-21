package registry

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	gAgents "github.com/totoual/gAgents/agent"
	pb "github.com/totoual/gAgents/protos/registry"
	chatgpt "github.com/totoual/gAgents/services/chatGPT"
	"google.golang.org/grpc"
)

const (
	SearchEvent          gAgents.EventType = "SearchEvent"
	TopicSuggestionEvent gAgents.EventType = "SuggestionEvent"
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
	// Emit event to ask GPT for topics
	responseChan := make(chan interface{})
	event := gAgents.Event{
		Type:         TopicSuggestionEvent,
		Payload:      a,
		ResponseChan: responseChan,
	}
	rs.eventDispatcher.Publish(event)

	// Wait for the response in the channel
	select {
	case response := <-responseChan:
		topics := response.([]string)
		// Successfully received a response
		fmt.Printf("Registry: we receivend the response from ChatGPT: %s", response)
		return &pb.RegistrationResponse{
			Success: true,
			Message: "Agent registered successfully.",
			Topics:  topics,
		}, nil
	case <-time.After(10 * time.Second): // 10 seconds timeout, adjust as needed
		// Timed out waiting for a response
		delete(rs.Agents, a.UniqueId)
		return &pb.RegistrationResponse{
			Success: false,
			Message: "Couldn't register!",
			Topics:  make([]string, 0),
		}, nil
	}

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

	return &pb.HealthCheckResponse{
		Success: true,
		Message: "Heartbeat received successfully",
	}, nil
}

func (rs *RegistryService) Search(ctx context.Context, sm *pb.SearchMessage) (*pb.SearchMessageResponse, error) {
	// For simplicity, just returning a success message without actually doing anything

	// Emit event to convert the Message to Kafka Message.
	responseChan := make(chan interface{})
	event := gAgents.Event{
		Type:         SearchEvent,
		Payload:      sm,
		ResponseChan: responseChan,
	}
	rs.eventDispatcher.Publish(event)

	select {
	case response := <-responseChan:

		sr := response.(chatgpt.SearchResult)
		// Build the kafka Message to publish
		// Successfully received a response
		fmt.Printf("Registry: we receivend the response from ChatGPT: %v", sr)
		fmt.Printf("The price is : %f", sr.PriceRange)
		return &pb.SearchMessageResponse{
			Success: true,
			Message: "Search request received.",
		}, nil
	case <-time.After(20 * time.Second): // 10 seconds timeout, adjust as needed
		// Timed out waiting for a response
		return &pb.SearchMessageResponse{
			Success: false,
			Message: "Search request timedout.",
		}, nil
	}
}

/*
Create functions to enable search functionality. The search functionaity will
create an event for the kafka.sendMessage function so we can inform the rest
of the agents about the search.
*/
