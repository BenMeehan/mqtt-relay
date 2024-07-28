package queue

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/benmeehan/mqtt-relay/config"
	"github.com/nats-io/nats.go"
)

// NATSQueue implements the Queue interface using NATS.
type NATSQueue struct {
	conn *nats.Conn
}

// NewNATSQueue creates a new NATSQueue instance with secure connection options.
func NewNATSQueue(cfg *config.Config) (*NATSQueue, error) {
	opts := []nats.Option{}

	if cfg.NATS.Username != "" && cfg.NATS.Password != "" {
		opts = append(opts, nats.UserInfo(cfg.NATS.Username, cfg.NATS.Password))
	}

	if cfg.NATS.CACertFile != "" {
		tlsConfig, err := newTLSConfig(cfg.NATS.CACertFile)
		if err != nil {
			return nil, err
		}
		opts = append(opts, nats.Secure(tlsConfig))
	}

	conn, err := nats.Connect(cfg.NATS.URL, opts...)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to NATS server", cfg.NATS.URL)
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

// newTLSConfig creates a TLS configuration for the NATS client.
func newTLSConfig(caCertFile string) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	return &tls.Config{
		RootCAs: caCertPool,
	}, nil
}
