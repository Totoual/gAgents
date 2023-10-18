package kafka

import (
	"log"

	"github.com/IBM/sarama"
	gAgents "github.com/totoual/gAgents/agent"
	pb "github.com/totoual/gAgents/protos/registry"
)

type KafkaConsumerService struct {
	brokers         []string
	consumer        sarama.Consumer
	eventDispatcher *gAgents.EventDispatcher
}

func NewKafkaConsumerService(brokers []string, ed *gAgents.EventDispatcher) (*KafkaConsumerService, error) {
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
		err := ks.ConsumeMessages(topic, ks.EmitEvent)
		if err != nil {
			log.Println("Failed to subscribe to topic", topic, ":", err)
		} else {
			log.Println("Successfully subscribed to topic", topic)
		}
	}

	event.ResponseChan <- "Subscribed to topics successfully."
}

func (ks *KafkaConsumerService) EmitEvent(message []byte) {
	event := gAgents.Event{
		Type:    "KafkaMessage",
		Payload: message,
	}

	ks.eventDispatcher.Publish(event)
}

func (ks *KafkaConsumerService) ConsumeMessages(topic string, handler func(message []byte)) error {

	consumer, err := ks.consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-consumer.Messages():
				handler(msg.Value)
			case err := <-consumer.Errors():
				// handle error
				log.Println(err)
			}
		}
	}()

	return nil
}
