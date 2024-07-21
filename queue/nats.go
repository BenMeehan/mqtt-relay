package queue

import (
	"log"

	"github.com/nats-io/nats.go"
)

const NATS_HOST = "nats://172.24.160.18:4222"

// NATSQueue implements the Queue interface using NATS.
type NATSQueue struct {
	conn *nats.Conn
}

// NewNATSQueue creates a new NATSQueue instance.
func NewNATSQueue() (*NATSQueue, error) {
	conn, err := nats.Connect(NATS_HOST)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to NATS host", NATS_HOST)
	return &NATSQueue{conn: conn}, nil
}

// Publish sends a message to a NATS topic.
func (q *NATSQueue) Publish(topic string, message []byte) error {
	if err := q.conn.Publish(topic, message); err != nil {
		log.Printf("Error publishing to NATS topic %s: %v", topic, err)
		return err
	}
	return nil
}

// Close disconnects from the NATS server.
func (q *NATSQueue) Close() error {
	q.conn.Close()
	return nil
}
