package mqtt

import (
	"log"

	"github.com/benmeehan/mqtt-relay/queue"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// DefaultMessageHandler implements the MessageHandler for processing incoming MQTT messages.
type DefaultMessageHandler struct {
	Queue         queue.Queue
	TopicMappings map[string]string
}

// HandleMessage processes incoming MQTT messages and publishes them to the queue.
func (h *DefaultMessageHandler) HandleMessage(client mqtt.Client, msg mqtt.Message) {
	queueTopic, ok := h.TopicMappings[msg.Topic()]
	if !ok {
		log.Printf("No queue topic mapping found for MQTT topic: %s", msg.Topic())
		return
	}

	log.Printf("Received message: %s from MQTT topic: %s\n", msg.Payload(), msg.Topic())

	if err := h.Queue.Publish(queueTopic, msg.Payload()); err != nil {
		log.Printf("Failed to publish to queue topic %s: %v", queueTopic, err)
	} else {
		log.Printf("Published message: %s to NATS topic: %s", msg.Payload(), queueTopic)
	}
}
