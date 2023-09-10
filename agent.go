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
	Receiver string
	Sender   string
	Type     string
	Content  string
}

type Server struct {
	pb.MessageServiceServer
	Agent *Agent
}

type MessageHandler interface {
	HandleMessage(message Message)
}

type Agent struct {
	name            string
	Addr            string
	InMessageQueue  chan Message
	OutMessageQueue chan Message
	grpcSrv         *grpc.Server
	messageHandlers map[string]MessageHandler
}

func NewAgent(name string, addr string) *Agent {
	return &Agent{
		name:            name,
		Addr:            addr,
		InMessageQueue:  make(chan Message),
		OutMessageQueue: make(chan Message),
		messageHandlers: make(map[string]MessageHandler),
	}
}

func (a *Agent) doSendMessage(m Message) (*pb.MessageResponse, error) {
	log.Printf("Trying to send the message to %v", m.Receiver)
	log.Printf("The content is: %v", m.Content)
	conn, err := grpc.Dial(m.Receiver, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Send the message via gRPC
	response, err := client.SendMessage(ctx, &pb.MessageRequest{
		Sender:   a.Addr,
		Receiver: m.Receiver,
		Content:  m.Content,
	})
	if err != nil {
		return nil, fmt.Errorf("Error sending message: %v", err)
	}

	return response, nil
}

func (s *Server) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Printf("SendMessage function was invoked with %v\n", in)

	message := Message{
		Sender:  in.Sender,
		Content: in.Content,
		Type:    in.Type,
	}
	s.Agent.InMessageQueue <- message
	log.Printf("Added a new message in the InMessageQueue: %v", message)
	return &pb.MessageResponse{
		Status: "OK",
	}, nil
}

func (a *Agent) ConsumeInMessages() {
	for {
		select {
		case message := <-a.InMessageQueue:
			a.DispatchMessage(message)
			// Process the received message here
			log.Printf("%s: Received message from %s: %s\n", a.name, message.Sender, message.Content)
		}
	}
}

func (a *Agent) ConsumeOutMessages() {
	for {
		log.Printf("Consuming messages")
		log.Printf("%v", len(a.OutMessageQueue))
		select {
		case message := <-a.OutMessageQueue:
			log.Printf("Consuming message: %v", message)
			a.doSendMessage(message)

			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (a *Agent) RegisterHandler(messageType string, handler MessageHandler) {
	a.messageHandlers[messageType] = handler
}

func (a *Agent) DispatchMessage(message Message) {
	handler, exists := a.messageHandlers[message.Type]
	if !exists {
		log.Printf("No handler found for message type: %s", message.Type)
		return
	}

	handler.HandleMessage(message)
}

func (a *Agent) Run(ctx context.Context) {
	// Start consuming messages
	go a.ConsumeInMessages()
	go a.ConsumeOutMessages()
	// Start the gRPC server.
	a.grpcSrv = grpc.NewServer()
	pb.RegisterMessageServiceServer(a.grpcSrv, &Server{Agent: a})
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
