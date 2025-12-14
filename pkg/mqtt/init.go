package mqtt

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"tj/config"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	Client paho.Client
}

func NewMQTTClient(clientId string) (*MQTTClient, error) {
	opts := paho.NewClientOptions()
	brokerURL := fmt.Sprintf("tcp://%s:%s", config.Cfg.MQTTBroker, config.Cfg.MQTTPort)

	opts.AddBroker(brokerURL)
	opts.SetClientID(clientId)
	opts.SetCleanSession(true)
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetKeepAlive(20 * time.Second)
	opts.SetPingTimeout(3 * time.Second)

	if config.Cfg.MQTTUseTLS == "true" {
		opts.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}

	opts.OnConnectionLost = func(c paho.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	}
	opts.OnConnect = func(c paho.Client) {
		log.Println("MQTT connected")
	}
	opts.OnReconnecting = func(c paho.Client, opts *paho.ClientOptions) {
		log.Println("MQTT reconnecting...")
	}

	client := paho.NewClient(opts)
	token := client.Connect()

	token.Wait()

	if token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	return &MQTTClient{Client: client}, nil
}

func (m *MQTTClient) Publish(topic string, payload []byte) error {
	if !m.Client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}

	token := m.Client.Publish(topic, 1, false, payload)

	token.Wait()

	return token.Error()
}

func (m *MQTTClient) Subscribe(topic string, handler paho.MessageHandler) error {
	if !m.Client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}

	token := m.Client.Subscribe(topic, 1, handler)

	token.Wait()

	return token.Error()
}

func (m *MQTTClient) Disconnect() {
	if m.Client.IsConnected() {
		m.Client.Disconnect(250)

		log.Println("MQTT disconnected")
	}
}
