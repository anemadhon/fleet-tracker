package controller

import (
	"encoding/json"
	"log"
	db "tj/pkg/database"
	geopkg "tj/pkg/geofence"
	model "tj/pkg/model"

	rmq "tj/pkg/rabbitmq"
)

type Worker struct {
	rmq *rmq.RabbitClient
	cfg rmq.RabbitConfig
}

func NewWorker(r *rmq.RabbitClient, cfg rmq.RabbitConfig) *Worker {
	return &Worker{rmq: r, cfg: cfg}
}

func (w *Worker) Start() error {
	msgs, err := rmq.ConsumeRMQWithConfig(w.rmq, w.cfg, true)
	if err != nil {
		return err
	}

	log.Printf("Geofence worker started, queue=%s routing=%s consumer=%s",
		w.cfg.QueueName, w.cfg.RoutingKey, w.cfg.ConsumerName)

	go func() {
		for d := range msgs {
			w.handleLocationMessage(d.Body)
		}
	}()

	return nil
}

func (w *Worker) handleLocationMessage(body []byte) {
	var loc model.VehicleLocation
	if err := json.Unmarshal(body, &loc); err != nil {
		log.Fatalf("invalid location payload: %v", err)
		return
	}

	var stations []model.BusStation
	if err := db.DB.Find(&stations).Error; err != nil {
		log.Fatalf("load bus_stations error: %v", err)
		return
	}

	const radius = 50.0
	for _, st := range stations {
		dist := geopkg.HaversineMeters(loc.Latitude, loc.Longitude, st.Latitude, st.Longitude)
		if dist <= radius {
			w.publishGeofenceEntry(&loc, &st, dist)
		}
	}
}

func (w *Worker) publishGeofenceEntry(loc *model.VehicleLocation, st *model.BusStation, dist float64) {
	evt := map[string]interface{}{
		"vehicle_id": loc.VehicleId,
		"event":      "geofence_entry",
		"location": map[string]interface{}{
			"latitude":  loc.Latitude,
			"longitude": loc.Longitude,
		},
		"timestamp": loc.Timestamp,
		"station": map[string]interface{}{
			"id":   st.Id,
			"name": st.Name,
		},
		"distance_m": dist,
	}
	b, err := json.Marshal(evt)
	if err != nil {
		log.Fatalf("marshal geofence_entry error: %v", err)
		return
	}
	if err := rmq.PublishRMQ(w.rmq, "fleet.events", "geofence.entry", b); err != nil {
		log.Fatalf("publish geofence.entry error: %v", err)
		return
	}

	log.Printf("geofence_entry vehicle=%s station=%s dist=%.1f m",
		loc.VehicleId, st.Name, dist)
}
