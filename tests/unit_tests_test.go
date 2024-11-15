package tests

import (
	"RTDS_API/services"
	"testing"
)

func TestInitKafkaProducer(t *testing.T) {
	streamID := "test-stream-id"
	err := services.InitKafkaProducer("localhost:9092", streamID)
	if err != nil {
		t.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
}

func TestProduceMessage(t *testing.T) {
	streamID := "test-stream-id"
	message := "test message"
	err := services.ProduceMessage(streamID, message)
	if err != nil {
		t.Errorf("Failed to produce message: %v", err)
	}
}
