package kafka

import (
	"log"
	"strings"

	"github.com/IBM/sarama"
	gAgents "github.com/totoual/gAgents/agent"
)

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

	return ks, nil
}

func (ks *KafkaProducerService) SendMessage(topic string, message []byte) error {
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
