package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"tj/config"
	"tj/pkg/mqtt"
	mock "tj/services/publisher/internal/controller"
)

func main() {
	config.Load()

	vehicleId := "B1234XYZ"
	clientId := "vehicle-" + vehicleId
	mqttClient, err := mqtt.NewMQTTClient(clientId)
	if err != nil {
		log.Fatalf("MQTT connect error: %v", err)
	}
	defer mqttClient.Disconnect()

	publisher := mock.NewPublisher(mqttClient.Client)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start publishing in goroutine
	go func() {
		publisher.PublishRandomMovement(vehicleId)
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down publisher...")
}
