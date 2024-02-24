package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
	gAgents "github.com/totoual/gAgents/agent"
)

const (
	KafkaSearchRequest = "KafkaSearchRequest"
)

type KafkaSearchMessage struct {
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

func NewKafkaSearchMessage(
	unique_id string,
	grpc_address string,
	object string,
	characteristics []string,
	category string,
	price_range float32,
	intended_use string,
	material_preferences []string,
	relevant_topics []string,
) (*KafkaSearchMessage, error) {
	return &KafkaSearchMessage{
		UniqueId:            unique_id,
		GrpcAddress:         grpc_address,
		Object:              object,
		Characteristics:     characteristics,
		Category:            category,
		PriceRange:          price_range,
		IntendedUse:         intended_use,
		MaterialPreferences: material_preferences,
		RelevantTopics:      relevant_topics,
	}, nil
}

type KafkaProducerService struct {
	brokers         []string
	producer        sarama.AsyncProducer
	eventDispatcher *gAgents.EventDispatcher
}

func NewKafkaProducerService(brokers []string, ed *gAgents.EventDispatcher, topics []string) (*KafkaProducerService, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	ks := &KafkaProducerService{
		brokers:         brokers,
		producer:        producer,
		eventDispatcher: ed,
	}
	ks.createTopics(topics)
	ks.eventDispatcher.Subscribe(KafkaSearchRequest, ks.handleSearchEvent)
	return ks, nil
}

func (ks *KafkaProducerService) handleSearchEvent(event gAgents.Event) {
	log.Printf("Handling a search Event!")
	search_request, ok := event.Payload.(*KafkaSearchMessage)
	if !ok {
		log.Println("Invalid payload type for Search Request event")
		return
	}
	jsonData, err := json.Marshal(search_request)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("Publishing a new Message to the Channel!")
	for _, topic := range search_request.RelevantTopics {
		ks.sendMessage(topic, jsonData)
	}
}

func (ks *KafkaProducerService) sendMessage(topic string, message []byte) error {
	log.Printf("sendMessage invoked! Sending the message!")
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	ks.producer.Input() <- msg
	return nil
}

func (ks *KafkaProducerService) createTopics(topics []string) error {
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

func (ks *KafkaProducerService) Close() error {
	if err := ks.producer.Close(); err != nil {
		return err
	}
	// Close other resources if needed
	return nil
}
