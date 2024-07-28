package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// MQTTConfig represents the MQTT configuration.
type MQTTConfig struct {
	Broker      string   `yaml:"broker"`
	ClientID    string   `yaml:"clientID"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	CACertFile  string   `yaml:"caCertFile"`
	ClientCertFile string `yaml:"clientCertFile"`
	ClientKeyFile   string `yaml:"clientKeyFile"`
	Topics      []TopicMapping `yaml:"topics"`
}

// TopicMapping maps an MQTT topic to a queue topic.
type TopicMapping struct {
	MQTTTopic  string `yaml:"mqttTopic"`
	QueueTopic string `yaml:"queueTopic"`
}

// Config represents the application's configuration.
type Config struct {
	MQTT MQTTConfig `yaml:"mqtt"`
	NATS struct {
		URL        string `yaml:"url"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		CACertFile string `yaml:"caCertFile"`
	} `yaml:"nats"`
}

// LoadConfig loads configuration from a YAML file.
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate validates the loaded configuration.
func (cfg *Config) Validate() error {
	if cfg.MQTT.Broker == "" || cfg.MQTT.ClientID == "" {
		return fmt.Errorf("MQTT broker and client ID must be provided")
	}
	if len(cfg.MQTT.Topics) == 0 {
		return fmt.Errorf("at least one MQTT topic must be configured")
	}
	return nil
}
