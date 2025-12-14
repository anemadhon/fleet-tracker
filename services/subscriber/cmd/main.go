package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"tj/config"
	db "tj/pkg/database"
	"tj/pkg/mqtt"
	rmq "tj/pkg/rabbitmq"
	sub "tj/services/subscriber/internal/controller"
)

func main() {
	config.Load()

	if err := db.Connect(); err != nil {
		log.Fatalf("Postgres init error: %v", err)
	}

	vehicleId := "B1234XYZ"
	clientId := "subcribe-vehicle-" + vehicleId
	mqttClient, err := mqtt.NewMQTTClient(clientId)
	if err != nil {
		log.Fatalf("MQTT connect error: %v", err)
	}
	defer mqttClient.Disconnect()

	rmqClient, err := rmq.Connect()
	if err != nil {
		log.Fatalf("RabbitMQ init error: %v", err)
	}
	defer rmqClient.Close()

	cfg := rmq.RabbitConfig{
		ExchangeName: "fleet.events",
		ExchangeType: "topic",
		QueueName:    "geofence_alerts",
		RoutingKey:   "location.raw",
	}

	if err := rmq.SetupRMQ(rmqClient, cfg); err != nil {
		log.Fatalf("RabbitMQ setup error: %v", err)
	}

	subscriber := sub.NewLocationSubscriber(rmqClient)
	topic := "/fleet/vehicle/+/location"
	if err := mqttClient.Subscribe(topic, subscriber.HandleMessage); err != nil {
		log.Fatalf("MQTT subscribe error: %v", err)
	}

	log.Printf("Subscriber listening on topic %s", topic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down subscriber...")
}
