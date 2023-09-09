package gAgents

import (
	"log"
	"net"

	pb "path_to_your_generated_grpc_code" // Import the generated gRPC package

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Agent struct {
	name    string
	grpcSrv *grpc.Server
}

func NewAgent(name string) *Agent {
	return &Agent{
		name: name,
	}
}

func (a *Agent) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	return &pb.MessageResponse{Status: "Message received"}, nil
}

func (a *Agent) ReceiveMessage(message Message) {
	log.Printf("%s: Received message from %s: %s\n", a.name, message.Sender, message.Content)
}

func (a *Agent) Run() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	a.grpcSrv = grpc.NewServer()
	pb.RegisterMessageServiceServer(a.grpcSrv, a)
	go a.grpcSrv.Serve(lis)
}
