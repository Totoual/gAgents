package kafka

import (
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
	gAgents "github.com/totoual/gAgents/agent"
	pb "github.com/totoual/gAgents/registry/proto"
	"github.com/totoual/gAgents/registry/services/registry"
)

type KafkaService struct {
	brokers         []string
	producer        sarama.AsyncProducer
	eventDispatcher *gAgents.EventDispatcher
}

func NewKafkaService(brokers []string, ed *gAgents.EventDispatcher, topics []string) (*KafkaService, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	ks := &KafkaService{
		brokers:         brokers,
		producer:        producer,
		eventDispatcher: ed,
	}
	ks.createTopics(topics)
	ed.Subscribe(registry.AgentRegisteredEventType, ks.handleRegistrationEvent)

	return ks, nil
}

func (ks *KafkaService) SendMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	ks.producer.Input() <- msg
	return nil
}

func (ks *KafkaService) ConsumeMessages(topic string, handler func(message []byte)) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewConsumer(ks.brokers, config)
	if err != nil {
		return err
	}

	consumer, err := client.ConsumePartition(topic, 0, sarama.OffsetOldest)
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

func (ks *KafkaService) handleRegistrationEvent(event gAgents.Event) {
	agentRegistration, ok := event.Payload.(*pb.AgentRegistration)
	if !ok {
		log.Println("Invalid payload type for AgentRegistered event")
		return
	}

	log.Println(agentRegistration)

	// TODO: Understand where we should register the user, based on
	// tags, capabilities and type.

	// TODO:// Subscribe the agent to the channel before you return.
	event.ResponseChan <- fmt.Sprintf("Topic created successfully.")
}

func (ks *KafkaService) createTopics(topics []string) error {
	// Find a better way to create the topic. This one is an example.

	adminConfig := sarama.NewConfig()
	adminConfig.Version = sarama.MaxVersion

	adminClient, err := sarama.NewClusterAdmin(ks.brokers, adminConfig)
	if err != nil {
		log.Println("Failed to create Kafka admin client:", err)
		return err
	}
	defer adminClient.Close()

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	for _, topic := range topics {
		err = adminClient.CreateTopic(topic, topicDetail, false)
		if err != nil {
			if strings.Contains(err.Error(), "Topic with this name already exists") {
				log.Println("Kafka Topic already exists. Returning the topic Name")
			} else {
				log.Println("Failed to create Kafka topic:", err)
				return err
			}
		}

		log.Println("Successfully created Kafka topic:", topic)
	}
	return nil
}

func (ks *KafkaService) Close() error {
	if err := ks.producer.Close(); err != nil {
		return err
	}
	// Close other resources if needed
	return nil
}
