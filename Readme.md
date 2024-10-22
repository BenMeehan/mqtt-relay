# MQTT Relay

## Overview

The MQTT Relay is a service that listens to MQTT topics from a broker and publishes the received messages to a queue system like NATS or Kafka. This project provides a modular and extensible solution for message relaying and transformation between different messaging systems.

It was written to solve an IoT use case like below, where there is a huge influx of messages (or metrics) being published through MQTT and the backend services may not be able to handle them all (due to database constraints). MQTT Relay acts as a queue source providing authentication and buffering the messages temporarily.

![MQTT Relay Architecture](_images/relay.png)

## Features

- **MQTT Client:** Connects to an MQTT broker and subscribes to specified topics.
- **Queue Integration:** Supports publishing received messages to NATS or Kafka queues.
- **Configurable:** Allows detailed configuration for MQTT topics and queue mappings.
- **Secure Connections:** Supports TLS/SSL for secure MQTT and NATS connections.
- **Logging:** Provides detailed logging for troubleshooting and monitoring.

## Getting Started

Follow these steps to set up and run the MQTT Relay.

### Prerequisites

- Go 1.18+ installed
- Docker (optional, for running NATS and Kafka in containers)
- Access to an MQTT broker
- Access to a NATS or Kafka broker (depending on your choice)

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/mqtt-relay.git
   cd mqtt-relay
   ```

2. **Install Dependencies**

   ```bash
   go mod tidy
   ```

### Configuration

Create a configuration file named `config/config.yaml`. Here’s a sample configuration:

```yaml
mqtt:
  broker: "tcp://broker.hivemq.com:1883"
  clientID: "mqtt_relay"
  username: "your_mqtt_username"
  password: "your_mqtt_password"
  caCertFile: "path/to/ca.crt"
  clientCertFile: "path/to/client.crt"
  clientKeyFile: "path/to/client.key"
  topics:
    - mqttTopic: "topic/one"
      queueTopic: "queue_topic_one"
    - mqttTopic: "topic/two"
      queueTopic: "queue_topic_two"
    - mqttTopic: "topic/three"
      queueTopic: "queue_topic_three"
nats:
  url: "tls://your_nats_server:4222"
  username: "your_nats_username"
  password: "your_nats_password"
  caCertFile: "path/to/nats/ca.crt"
```

- **mqtt.broker:** The MQTT broker URL.
- **mqtt.clientID:** The client ID for the MQTT connection.
- **mqtt.username:** The username for the MQTT broker.
- **mqtt.password:** The password for the MQTT broker.
- **mqtt.caCertFile:** The path to the CA certificate for MQTT.
- **mqtt.clientCertFile:** The path to the client certificate for MQTT.
- **mqtt.clientKeyFile:** The path to the client key for MQTT.
- **mqtt.topics:** Mapping between MQTT topics and queue topics.
- **nats.url:** The NATS broker URL.
- **nats.username:** The username for the NATS broker.
- **nats.password:** The password for the NATS broker.
- **nats.caCertFile:** The path to the CA certificate for NATS.

### Running the Project

1. **Build the Application**

   ```bash
   go build -o mqtt-relay
   ```

2. **Run the Application**

   ```bash
   ./mqtt-relay
   ```

   The application will start, connect to the MQTT broker, subscribe to the specified topics, and publish received messages to the configured queue system.

### Using Docker (Optional)

You can run NATS and Kafka using Docker if you don’t have them set up.

**Run NATS**

```bash
docker run -d --name nats -p 4222:4222 nats:latest
```

**Run Kafka**

```bash
docker run -d --name kafka -p 9092:9092 -e KAFKA_ADVERTISED_LISTENERS=INSIDE://kafka:9092,OUTSIDE://localhost:9092 -e KAFKA_LISTENER_NAME=INSIDE -e KAFKA_LISTENERS=INSIDE://0.0.0.0:9092,OUTSIDE://0.0.0.0:9092 wurstmeister/kafka:latest
```

### Testing

To test the application, ensure that both the MQTT broker and the queue system (NATS or Kafka) are running. The application should start receiving messages from the MQTT broker and publishing them to the queue.

### Logging and Debugging

Logs are output to the standard console. You can add or modify logging in the source code to adjust verbosity or format.

### Contributing

1. **Fork the Repository**
2. **Create a New Branch**
3. **Commit Your Changes**
4. **Push to Your Branch**
5. **Open a Pull Request**

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### Contact

For any issues or questions, please open an issue on GitHub or contact [benmeehan111@gmail.com](mailto:benmeehan111@gmail.com).
