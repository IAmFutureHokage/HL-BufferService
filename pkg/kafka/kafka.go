package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaConfig struct {
	BrokerList []string `mapstructure:"broker_list"`
	Topic      string   `mapstructure:"topic"`
}

func NewKafkaProducer(config KafkaConfig) (sarama.SyncProducer, error) {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForLocal       // Принимать подтверждение после записи в локальный лог
	producerConfig.Producer.Compression = sarama.CompressionSnappy   // Используем сжатие Snappy
	producerConfig.Producer.Flush.Frequency = 500 * time.Millisecond // Как часто отправлять накопленные сообщения
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.BrokerList, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	return producer, nil
}

func SendMessageToKafka(producer sarama.SyncProducer, topic string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	return nil
}
