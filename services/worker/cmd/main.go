package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"tj/config"
	db "tj/pkg/database"
	rmq "tj/pkg/rabbitmq"
	geo "tj/services/worker/internal/controller"
)

func main() {
	config.Load()

	if err := db.Connect(); err != nil {
		log.Fatalf("Postgres init error: %v", err)
	}

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
		ConsumerName: "worker-geofence",
	}

	if err := rmq.SetupRMQ(rmqClient, cfg); err != nil {
		log.Fatalf("RabbitMQ setup error: %v", err)
	}

	worker := geo.NewWorker(rmqClient, cfg)
	if err := worker.Start(); err != nil {
		log.Fatalf("worker start error: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down geofence worker")
}
