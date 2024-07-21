package queue

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaQueue represents a Kafka queue.
type KafkaQueue struct {
	writer *kafka.Writer
}

// NewKafkaQueue creates a new KafkaQueue instance with the given broker and topic.
func NewKafkaQueue(broker string, topic string) (*KafkaQueue, error) {
	writer := &kafka.Writer{
		Addr:  kafka.TCP(broker),
		Topic: topic,
	}

	// Test connection by sending a test message
	if err := testKafkaConnection(writer); err != nil {
		return nil, err
	}

	log.Printf("Connected to Kafka broker %s on topic %s", broker, topic)
	return &KafkaQueue{writer: writer}, nil
}

// testKafkaConnection sends a test message to verify the connection.
func testKafkaConnection(writer *kafka.Writer) error {
	testMsg := kafka.Message{
		Key:   []byte("test-key"),
		Value: []byte("test-value"),
	}

	err := writer.WriteMessages(context.Background(), testMsg)
	if err != nil {
		log.Printf("Failed to send test message: %v", err)
		return err
	}

	log.Println("Kafka connection test message sent successfully.")
	return nil
}

// Publish sends a message to the specified Kafka topic.
func (q *KafkaQueue) Publish(topic string, message []byte) error {
	err := q.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Value: message,
	})
	if err != nil {
		log.Printf("Error publishing message to topic %s: %v", topic, err)
		return err
	}
	return nil
}

// Close closes the Kafka writer.
func (q *KafkaQueue) Close() error {
	if err := q.writer.Close(); err != nil {
		log.Printf("Error closing Kafka writer: %v", err)
		return err
	}
	log.Println("Kafka writer closed successfully.")
	return nil
}
