package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	gAgents "github.com/totoual/gAgents/agent"
	pb "github.com/totoual/gAgents/protos/registry"
)

const (
	KafkaRegistration gAgents.EventType = "KafkaRegistrationMessage"
)

type KafkaConsumerService struct {
	brokers         []string
	consumer        sarama.Consumer
	eventDispatcher *gAgents.EventDispatcher
}

type KafkaConsumerMessage struct {
	UniqueId            string   `json:"unique_id"`
	GrpcAddress         string   `json:"grpc_address"`
	Object              string   `json:"object"`
	Characteristics     []string `json:"characteristics"`
	Category            string   `json:"category"`
	PriceRange          float32  `json:"price_range"`
	IntendedUse         string   `json:"intended_use"`
	MaterialPreferences []string `json:"material_preferences"`
	RelevantTopics      []string `json:"relevant_topics"`
}

func NewKafkaConsumerService(brokers []string, ed *gAgents.EventDispatcher, e gAgents.EventType) (*KafkaConsumerService, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	ks := &KafkaConsumerService{
		brokers:         brokers,
		consumer:        consumer,
		eventDispatcher: ed,
	}

	ed.Subscribe(e, ks.handleRegistrationEvent)

	return ks, nil
}

func (ks *KafkaConsumerService) handleRegistrationEvent(event gAgents.Event) {
	registration_response, ok := event.Payload.(*pb.RegistrationResponse)
	if !ok {
		log.Println("Invalid payload type for RegistrationResponse event")
		return
	}

	log.Println(registration_response)

	for _, topic := range registration_response.Topics {
		err := ks.ConsumeMessages(topic, KafkaRegistration, ks.EmitEvent)
		if err != nil {
			log.Println("Failed to subscribe to topic", topic, ":", err)
		} else {
			log.Println("Successfully subscribed to topic", topic)
		}
	}

	event.ResponseChan <- "Subscribed to topics successfully."
}

func (ks *KafkaConsumerService) EmitEvent(kafka_type gAgents.EventType, message []byte) {
	fmt.Printf("Publishing a new event needs to be handled!")
	fmt.Println(kafka_type)
	var msg KafkaConsumerMessage
	err := json.Unmarshal(message, &msg)
	if err != nil {
		fmt.Errorf("Could not unmarsal the message.")
	}
	fmt.Println("The message I received is: %s", msg.Category)
	event := gAgents.Event{
		Type:    kafka_type,
		Payload: message,
	}

	ks.eventDispatcher.Publish(event)
}

func (ks *KafkaConsumerService) ConsumeMessages(topic string, kafka_type gAgents.EventType, handler func(gAgents.EventType, []byte)) error {

	consumer, err := ks.consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-consumer.Messages():
				handler(kafka_type, msg.Value)
			case err := <-consumer.Errors():
				// handle error
				log.Println(err)
			}
		}
	}()

	return nil
}
