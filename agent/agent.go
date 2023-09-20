package gAgents

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	pb "github.com/totoual/gAgents/generated/proto" // Import the generated gRPC package

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Envelope struct {
	Receiver string
	Sender   string
	Protocol string
	Message  []byte
}

func NewEnvelope(m Message) *Envelope {
	sm, err := m.Serialize()
	if err != nil {
		log.Panic("Cannot serialise this message!")
	}
	return &Envelope{
		Sender:   m.GetSender(),
		Receiver: m.GetReceiver(),
		Protocol: m.GetProtocol(),
		Message:  sm,
	}
}

func (e Envelope) ToMessage(m Message) (Message, error) {
	err := json.Unmarshal(e.Message, &m)
	if err != nil {
		return nil, fmt.Errorf("error deserializing message: %v", err)
	}
	return m, nil
}

type Message interface {
	GetProtocol() string
	GetPerformative() int
	GetReceiver() string
	GetSender() string
	Serialize() ([]byte, error)
}

type Server struct {
	pb.MessageServiceServer
	Agent *Agent
}

type Handler interface {
	HandleMessage(envelope Envelope)
}

type Act interface {
	Perform(a *Agent)
	GetInterval() time.Duration
}

type BusinessLogic interface {
	Apply() Message
}

type Agent struct {
	name            string
	Addr            string
	InMessageQueue  chan Envelope
	OutMessageQueue chan Envelope
	grpcSrv         *grpc.Server
	handlers        map[string]Handler
	acts            []Act
	BusinessLogic   map[string]BusinessLogic
	TaskScheduler   *TaskScheduler
	ctx             context.Context
	Cancel          context.CancelFunc
}

func NewAgent(name string, addr string) *Agent {
	ctx, cancel := context.WithCancel(context.Background())
	return &Agent{
		name:            name,
		Addr:            addr,
		InMessageQueue:  make(chan Envelope),
		OutMessageQueue: make(chan Envelope),
		handlers:        make(map[string]Handler),
		acts:            make([]Act, 0),
		BusinessLogic:   make(map[string]BusinessLogic),
		TaskScheduler:   NewTaskScheduler(ctx),
		ctx:             ctx,
		Cancel:          cancel,
	}
}

func (a *Agent) doSendEnvelope(e Envelope) (*pb.MessageResponse, error) {
	log.Printf("Sending the message to %v", e.Receiver)
	conn, err := grpc.Dial(e.Receiver, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Send the message via gRPC
	response, err := client.SendMessage(ctx, &pb.MessageRequest{
		Sender:   a.Addr,
		Receiver: e.Receiver,
		Message:  e.Message,
		Uuid:     uuid.New().String(),
		Protocol: e.Protocol,
	})
	if err != nil {
		return nil, fmt.Errorf("error sending message: %v", err)
	}

	return response, nil
}

func (s *Server) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Printf("SendMessage function was invoked with %v\n", in)

	envelope := Envelope{
		Sender:   in.Sender,
		Message:  in.Message,
		Protocol: in.Protocol,
	}
	s.Agent.InMessageQueue <- envelope
	log.Printf("Added a new message in the InMessageQueue")
	return &pb.MessageResponse{
		Status: "OK",
	}, nil
}

func (a *Agent) ConsumeInMessages() {
	for {
		select {
		case envelope := <-a.InMessageQueue:
			a.DispatchMessage(envelope)
			// Process the received message here
			log.Printf("%s: Received message from %s\n", a.name, envelope.Sender)
		case <-a.ctx.Done():
			// Agent's context has been canceled, terminate the goroutine
			return
		}
	}
}

func (a *Agent) ConsumeOutMessages() {
	for {
		log.Printf("Consuming messages")
		log.Printf("%v", len(a.OutMessageQueue))
		select {
		case envelope := <-a.OutMessageQueue:
			log.Printf("Consuming message: %v", envelope)
			a.doSendEnvelope(envelope)

			time.Sleep(time.Millisecond * 100)
		case <-a.ctx.Done():
			// Agent's context has been canceled, terminate the goroutine
			return
		}
	}
}

func (a *Agent) DispatchMessage(envelope Envelope) {
	handler, exists := a.handlers[envelope.Protocol]
	if !exists {
		log.Printf("No handler found for message type: %s", envelope.Protocol)
		return
	}

	handler.HandleMessage(envelope)
}

func (a *Agent) RegisterHandler(messageType string, handler Handler) {
	a.handlers[messageType] = handler
}

func (a *Agent) RegisterAct(act Act) {
	a.acts = append(a.acts, act)
}

func (a *Agent) RegisterBusinessLogic(messageType string, logic BusinessLogic) {
	a.BusinessLogic[messageType] = logic
}

func (a *Agent) PerformActs() {
	// Perform initial setup acts
	// This could include sending messages, performing tasks, etc.

	// Loop through all registered acts
	for _, act := range a.acts {
		interval := act.GetInterval()

		// If act is periodic, create a ticker to control the interval
		if interval > 0 {
			ticker := time.NewTicker(interval)
			go func(act Act, ctx context.Context) {
				for {
					select {
					case <-ticker.C:
						act.Perform(a)
					case <-ctx.Done():
						return
					}
				}
			}(act, a.ctx)
		} else {
			// If act is not periodic, perform it immediately
			act.Perform(a)
		}
	}
}

func (a *Agent) Run() {
	// Start consuming messages
	go a.ConsumeInMessages()
	go a.ConsumeOutMessages()
	go a.PerformActs()
	go a.TaskScheduler.ExecuteTasks()
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
	case <-a.ctx.Done():
		// Context has been canceled, stop the server gracefully
		a.grpcSrv.Stop()
	}
}
