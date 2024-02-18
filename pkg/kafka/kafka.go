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

type MessageProducer interface {
	Serialize() ([]byte, error)
}

func NewKafkaProducer(config KafkaConfig) (sarama.SyncProducer, error) {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForLocal      // Принимать подтверждение после записи в локальный лог
	producerConfig.Producer.Compression = sarama.CompressionSnappy  // Используем сжатие Snappy
	producerConfig.Producer.Flush.Frequency = 50 * time.Millisecond // Как часто отправлять накопленные сообщения
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.BrokerList, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать кафка-продюссера: %v", err)
	}
	return producer, nil
}

func SendMessageToKafka(producer sarama.SyncProducer, topic string, messageProducer MessageProducer) error {
	messageBytes, err := messageProducer.Serialize()
	if err != nil {
		return fmt.Errorf("Ошибка серилизации: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(messageBytes),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("Не удалось отправить сообщение в Kafka: %v", err)
	}

	return nil
}
