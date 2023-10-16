package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaService struct {
	brokers  []string
	producer sarama.AsyncProducer
}

func NewKafkaService(brokers []string) (*KafkaService, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaService{
		brokers:  brokers,
		producer: producer,
	}, nil
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

func (ks *KafkaService) Close() error {
	if err := ks.producer.Close(); err != nil {
		return err
	}
	// Close other resources if needed
	return nil
}
