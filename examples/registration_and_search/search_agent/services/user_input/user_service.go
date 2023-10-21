package userinput

import (
	"context"
	"fmt"

	pb "github.com/totoual/gAgents/examples/registration_and_search/search_agent/services/user_input/proto"
)

type UserInteractionService struct {
	pb.UnimplementedUserInteractionServer
}

func (s *UserInteractionService) UserSearch(ctx context.Context, in *pb.UserSearchMessage) (*pb.UserSearchMessageResponse, error) {
	// Implement your search logic here...
	// For example:
	fmt.Printf("Received search request from %v: %v\n", in.UniqueId, in.Description)

	// Assuming search was successful:
	return &pb.UserSearchMessageResponse{
		Success: true,
		Message: "Search completed successfully.",
	}, nil
}
