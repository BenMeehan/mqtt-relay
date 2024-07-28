package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/benmeehan/mqtt-relay/config"
	"github.com/benmeehan/mqtt-relay/queue"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

// MQTTClient manages an MQTT client and its configuration.
type MQTTClient struct {
	Client         mqtt.Client
	Config         *config.Config
	Queue          queue.Queue
	MessageHandler mqtt.MessageHandler
}

// NewMQTTClient creates a new MQTT client with the provided configuration and message handler.
func NewMQTTClient(cfg *config.Config, q queue.Queue, handler mqtt.MessageHandler) (*MQTTClient, error) {
	clientID := fmt.Sprintf("%s-%s", cfg.MQTT.ClientID, uuid.New().String())

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTT.Broker)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(handler)

	if cfg.MQTT.Username != "" && cfg.MQTT.Password != "" {
		opts.SetUsername(cfg.MQTT.Username)
		opts.SetPassword(cfg.MQTT.Password)
	}

	if cfg.MQTT.CACertFile != "" && cfg.MQTT.ClientCertFile != "" && cfg.MQTT.ClientKeyFile != "" {
		tlsConfig, err := newTLSConfig(cfg.MQTT.CACertFile, cfg.MQTT.ClientCertFile, cfg.MQTT.ClientKeyFile)
		if err != nil {
			return nil, err
		}
		opts.SetTLSConfig(tlsConfig)
	}

	client := mqtt.NewClient(opts)

	return &MQTTClient{
		Client:        client,
		Config:        cfg,
		Queue:         q,
		MessageHandler: handler,
	}, nil
}

// Connect connects the MQTT client to the broker.
func (c *MQTTClient) Connect() error {
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Error connecting to broker: %v", token.Error())
		return token.Error()
	}
	return nil
}

// Subscribe subscribes to the configured topics.
func (c *MQTTClient) Subscribe() error {
	for _, mapping := range c.Config.MQTT.Topics {
		if token := c.Client.Subscribe(mapping.MQTTTopic, 2, nil); token.Wait() && token.Error() != nil {
			log.Printf("Error subscribing to topic %s: %v", mapping.MQTTTopic, token.Error())
			return fmt.Errorf("failed to subscribe to topic %s: %v", mapping.MQTTTopic, token.Error())
		}
		log.Printf("Subscribed to MQTT topic: %s", mapping.MQTTTopic)
	}
	return nil
}

// Disconnect disconnects the MQTT client from the broker.
func (c *MQTTClient) Disconnect() {
	c.Client.Disconnect(250)
}

// newTLSConfig creates a TLS configuration for the MQTT client.
func newTLSConfig(caCertFile, clientCertFile, clientKeyFile string) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate and key: %w", err)
	}

	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{clientCert},
	}, nil
}
