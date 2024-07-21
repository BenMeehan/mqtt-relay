package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/benmeehan/mqtt-relay/config"
	"github.com/benmeehan/mqtt-relay/mqtt"
	"github.com/benmeehan/mqtt-relay/queue"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create NATS queue
	q, err := queue.NewNATSQueue()
	if err != nil {
		log.Fatalf("Failed to create NATS queue: %v", err)
	}
	defer q.Close()

	// Create topic mapping
	topicMappings := make(map[string]string)
	for _, mapping := range cfg.MQTT.Topics {
		topicMappings[mapping.MQTTTopic] = mapping.QueueTopic
	}

	// Initialize MQTT message handler
	handler := &mqtt.DefaultMessageHandler{Queue: q, TopicMappings: topicMappings}

	// Create and configure MQTT client
	client := mqtt.NewMQTTClient(cfg, q, handler.HandleMessage)

	// Connect to the MQTT broker
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
	}
	log.Println("Connected to broker with config", client.Config)

	// Subscribe to topics
	if err := client.Subscribe(); err != nil {
		log.Fatalf("Failed to subscribe to topics: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the client
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc

	client.Disconnect()
	log.Println("Disconnected from broker")
}
