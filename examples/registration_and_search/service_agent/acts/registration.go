package acts

import (
	"context"
	"fmt"
	"log"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/registration_and_search/service_agent/config"
	pb "github.com/totoual/gAgents/protos/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const AgentRegisteredEventType = "AgentRegistered"

type RegistrationAct struct {
	Config config.AgentConfig
	Event  string
	ed     *gAgents.EventDispatcher
}

func NewRegistrationAct(config config.AgentConfig, ed *gAgents.EventDispatcher) *RegistrationAct {
	return &RegistrationAct{
		Config: config,
		ed:     ed,
		Event:  AgentRegisteredEventType,
	}
}

// Perform is the method called when the act is performed.
func (r *RegistrationAct) Perform(a *gAgents.Agent) {
	// Perform action...
	fmt.Println(" \n RegistrationAct performed.")
	fmt.Println(r.Config)
	log.Printf("Sending the message to %v", r.Config.RegistryURL)
	conn, err := grpc.Dial(r.Config.RegistryURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewAgentRegistryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Send the message via gRPC
	response, err := client.RegisterAgent(ctx, &pb.AgentRegistration{
		UniqueId:     r.Config.UniqueID,
		GrpcAddress:  r.Config.AgentURL,
		AgentType:    r.Config.AgentType,
		Capabilities: r.Config.Capabilities,
		Metadata: &pb.Metadata{
			SoftwareVersion: r.Config.Metadata.SoftwareVersion,
			Location: &pb.Location{
				Latitude:  r.Config.Metadata.Location.Latitude,
				Longitude: r.Config.Metadata.Location.Longitude,
			},
			Region:       r.Config.Metadata.Region,
			Organization: r.Config.Metadata.Organization,
		},
		Status: &pb.Status{
			CurrentStatus: pb.Status_AgentStatus(pb.Status_AgentStatus_value[r.Config.Status.CurrentStatus]),
		},
		AuthData: &pb.Authentication{
			Token:     r.Config.AuthData.Token,
			PublicKey: r.Config.AuthData.PublicKey,
		},
		ContactInfo: &pb.ContactInformation{
			Email:            r.Config.ContactInfo.Email,
			SecondaryChannel: r.Config.ContactInfo.SecondaryChannel,
		},
		Tags: r.Config.Tags,
	})
	if err != nil {
		fmt.Errorf("error sending message: %v", err)
	} else {

		// Emit an event that we managed to register. So we can subscribe to topics
		fmt.Println(response.Message)
		responseChan := make(chan interface{})
		event := gAgents.Event{
			Type:         AgentRegisteredEventType,
			Payload:      response,
			ResponseChan: responseChan,
		}
		r.ed.Publish(event)

		select {
		case response := <-responseChan:
			// Successfully received a response
			fmt.Printf("Subscribed to the topics: %s", response)
		case <-time.After(10 * time.Second): // 10 seconds timeout, adjust as needed
			// Timed out waiting for a response
			fmt.Printf("Failed to subscribe to any topic!")
		}
	}

}

// GetInterval returns the interval at which the act should be performed.
func (a *RegistrationAct) GetInterval() time.Duration {
	return 0 // Perform every 2 seconds
}
