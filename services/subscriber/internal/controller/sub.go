package controller

import (
	"encoding/json"
	"log"
	"sync"

	db "tj/pkg/database"
	rmq "tj/pkg/rabbitmq"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	model "tj/pkg/model"
)

type LocationSubscriber struct {
	mu      sync.Mutex
	buffer  []model.MQTTLocationStruct
	maxSize int
	rmq     *rmq.RabbitClient
}

func NewLocationSubscriber(rmq *rmq.RabbitClient) *LocationSubscriber {
	return &LocationSubscriber{
		buffer:  make([]model.MQTTLocationStruct, 0, 8),
		maxSize: 8,
		rmq:     rmq,
	}
}

func (h *LocationSubscriber) HandleMessage(client mqtt.Client, msg mqtt.Message) {
	var loc model.MQTTLocationStruct

	if err := json.Unmarshal(msg.Payload(), &loc); err != nil {
		log.Fatalf("invalid payload on topic %s: %v", msg.Topic(), err)
		return
	}
	if loc.VehicleId == "" {
		log.Fatalf("missing vehicle_id on topic %s", msg.Topic())
		return
	}

	record := model.MQTTLocationStruct{
		VehicleId: loc.VehicleId,
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
		Timestamp: loc.Timestamp,
	}

	if err := db.DB.Create(&record).Error; err != nil {
		log.Fatalf("GORM insert error: %v", err)
		return
	}
	// h.mu.Lock()

	// h.buffer = append(h.buffer, record)

	// if len(h.buffer) >= h.maxSize {
	// 	toInsert := h.buffer
	// 	h.buffer = make([]model.MQTTLocationStruct, 0, h.maxSize)
	// 	h.mu.Unlock()

	// 	if err := db.DB.Create(&toInsert).Error; err != nil {
	// 		log.Fatalf("GORM batch insert error: %v", err)
	// 	} else {
	// 		log.Printf("batch inserted %d records", len(toInsert))
	// 	}

	// 	return
	// }

	// h.mu.Unlock()

	log.Printf("stored (gorm) vehicle=%s lat=%.6f lon=%.6f ts=%d",
		record.VehicleId, record.Latitude, record.Longitude, record.Timestamp)

	paylaodInBytes, err := json.Marshal(record)
	if err != nil {
		log.Fatalf("err parse payload to bytes: %v", err)
	}
	if err = rmq.PublishRMQ(h.rmq, "fleet.events", "location.raw", paylaodInBytes); err != nil {
		log.Fatalf("publish location.raw error: %v", err)
	}
}
