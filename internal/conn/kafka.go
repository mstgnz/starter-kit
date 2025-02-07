package conn

import (
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

type Kafka struct {
	sarama.SyncProducer
}

func (k *Kafka) ConnectKafka() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
	url := os.Getenv("BROKER_URL")
	if len(url) != 0 {
		producer, err := sarama.NewSyncProducer([]string{url}, config)
		if err != nil {
			log.Println("Failed Kafka Connection", err.Error())
		} else {
			log.Println("Kafka Connected")
		}
		k.SyncProducer = producer
	} else {
		log.Println("BROKER_URL Not Found!")
	}
}

func (k *Kafka) CloseKafka() {
	if err := k.Close(); err != nil {
		log.Println("Failed to close Kafka Producer:", err.Error())
	} else {
		log.Println("Kafka Producer Closed")
	}
}

func (k *Kafka) PushCommentToQueue(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := k.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
