package services

import (
	"RTDS_API/config"
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var KafkaWriters = make(map[string]*kafka.Writer)

func InitKafkaProducer(brokerAddress, streamID string) error {
	if _, exists := KafkaWriters[streamID]; !exists {
		mechanism, err := scram.Mechanism(scram.SHA256, "Nik0maybe", "FJ57UxjIQFvdQRDNbH5XPQyqKEHMBr")
		if err != nil {
			return fmt.Errorf("failed to create SASL mechanism: %w", err)
		}

		writer := &kafka.Writer{
			Addr:     kafka.TCP(brokerAddress),
			Topic:    "stream-topic",
			Balancer: &kafka.LeastBytes{},
			Transport: &kafka.Transport{
				TLS:  &tls.Config{},
				SASL: mechanism,
			},
		}

		KafkaWriters[streamID] = writer
		log.Println("Kafka producer initialized for stream:", streamID)
	}
	return nil
}

func ProduceMessage(streamID, message string) error {
	writer, exists := KafkaWriters[streamID]
	if !exists {
		return fmt.Errorf("kafka writer for stream ID %s not initialized", streamID)
	}

	log.Printf("[Stream %s] Producing message to Kafka: %s", streamID, message)
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(streamID),
		Value: []byte(message),
	})
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}
	return nil
}

func ConsumeMessages(streamID string, resultsChan chan<- string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{config.BrokerAddress},
		Topic:       "stream-topic",
		GroupID:     "consumer-group-" + streamID,
		StartOffset: kafka.LastOffset,
	})

	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("[Stream %s] Error closing Kafka reader: %v", streamID, err)
		}
		close(resultsChan)
	}()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[Stream %s] Error reading message from Kafka: %v", streamID, err)
			continue
		}

		processedData := "Processed: " + string(msg.Value)
		log.Printf("[Stream %s] Consumed and processed message: %s", streamID, processedData)
		resultsChan <- processedData
	}
}
