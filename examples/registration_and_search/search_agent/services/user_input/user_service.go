package userinput

import (
	"context"
	"fmt"
	"time"

	gAgents "github.com/totoual/gAgents/agent"
	"github.com/totoual/gAgents/examples/registration_and_search/search_agent/config"
	pb "github.com/totoual/gAgents/examples/registration_and_search/search_agent/services/user_input/proto"
	rb "github.com/totoual/gAgents/protos/registry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserInteractionService struct {
	pb.UnimplementedUserInteractionServer
	eventDispatcher *gAgents.EventDispatcher
	Config          config.AgentConfig
}

func NewUserInteractionService(ed *gAgents.EventDispatcher, cfg config.AgentConfig) *UserInteractionService {
	return &UserInteractionService{
		eventDispatcher: ed,
		Config:          cfg,
	}
}

func (us *UserInteractionService) Init(srv *grpc.Server) {
	pb.RegisterUserInteractionServer(srv, us)
}

func (s *UserInteractionService) UserSearch(ctx context.Context, in *pb.UserSearchMessage) (*pb.UserSearchMessageResponse, error) {
	// Implement your search logic here...
	// For example:
	fmt.Printf("Received search request from %v: %v\n", in.UniqueId, in.Description)

	conn, err := grpc.Dial(s.Config.RegistryURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := rb.NewAgentRegistryClient(conn)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		response, err := client.Search(ctx, &rb.SearchMessage{
			UniqueId:    s.Config.UniqueID,
			GrpcAddress: s.Config.AgentURL,
			Description: in.Description,
		})

		if err != nil {
			fmt.Errorf("error sending message: %v", err)
		} else {
			// Handle the search response...
			fmt.Printf("Search response: %v\n", response.Message)
		}
	}()

	// Immediately respond to the user indicating the search has been initiated.
	return &pb.UserSearchMessageResponse{
		Success: true,
		Message: "Search request received and is being processed.",
	}, nil
}
