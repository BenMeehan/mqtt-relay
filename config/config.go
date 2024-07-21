package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// MQTTConfig represents the MQTT configuration.
type MQTTConfig struct {
	Broker    string `yaml:"broker"`
	ClientID  string `yaml:"clientID"`
	Topics    []TopicMapping `yaml:"topics"`
}

// TopicMapping maps an MQTT topic to a queue topic.
type TopicMapping struct {
	MQTTTopic string `yaml:"mqttTopic"`
	QueueTopic string `yaml:"queueTopic"`
}

// Config represents the application's configuration.
type Config struct {
	MQTT MQTTConfig `yaml:"mqtt"`
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

	return &config, nil
}
