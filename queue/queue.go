package queue

// Queue interface to allow easy addition of other message queues in the future
type Queue interface {
	Publish(topic string, message []byte) error
	Close() error
}
