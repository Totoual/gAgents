package gAgents

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/totoual/gAgents/generated/proto" // Import the generated gRPC package

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Message struct {
	Sender  string
	Content string
}

type Agent struct {
	name            string
	Addr            string
	InMessageQueue  chan Message
	OutMessageQueue chan Message
	grpcSrv         *grpc.Server
}

type Server struct {
	pb.MessageServiceServer
}

func NewAgent(name string, addr string) *Agent {
	return &Agent{
		name:            name,
		Addr:            addr,
		InMessageQueue:  make(chan Message),
		OutMessageQueue: make(chan Message),
	}
}

func (a *Agent) DoSendMessage(addr string, receiver string, content string) (*pb.MessageResponse, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Send the message via gRPC
	response, err := client.SendMessage(ctx, &pb.MessageRequest{
		Sender:   a.name,
		Receiver: receiver,
		Content:  content,
	})
	if err != nil {
		return nil, fmt.Errorf("Error sending message: %v", err)
	}

	return response, nil
}

func (a *Server) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Printf("SendMessage function was invoked with %v\n", in)
	return &pb.MessageResponse{
		Status: "OK",
	}, nil
}

func (a *Agent) SendAsyncMessage(receiver string, content string) {
	a.OutMessageQueue <- Message{
		Sender:  a.name,
		Content: content,
	}
}

func (a *Agent) HandleReceivedMessage(message Message) {
	// Process the received message here or pass it to a handler.
	log.Printf("%s: Received message from %s: %s\n", a.name, message.Sender, message.Content)
}

func (a *Agent) Run(ctx context.Context) {
	a.grpcSrv = grpc.NewServer()
	pb.RegisterMessageServiceServer(a.grpcSrv, &Server{})
	lis, err := net.Listen("tcp", a.Addr)

	if err != nil {
		log.Fatalf("Failed to list on: %v\n", err)
	}
	log.Printf("Listening on %s\n", a.Addr)

	// Start the server
	if err := a.grpcSrv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	select {
	case <-ctx.Done():
		// Context has been canceled, stop the server gracefully
		a.grpcSrv.Stop()
	}
}
