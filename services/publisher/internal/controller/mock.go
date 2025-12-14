package mock

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	model "tj/pkg/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MockPublisher struct {
	client mqtt.Client
	lat    float64
	lon    float64
}

func NewPublisher(client mqtt.Client) *MockPublisher {
	// Initialize random seed
	rand.Seed(time.Now().Unix())

	// Set posisi awal (Jakarta Office area)
	return &MockPublisher{
		client: client,
		lat:    -6.2088,
		lon:    106.8456,
	}
}

// PublishRandomMovement publishes location data every 2 seconds with random movement
func (p *MockPublisher) PublishRandomMovement(vehicleId string) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	log.Printf("Starting mock publisher for vehicle: %s\n", vehicleId)
	log.Printf("Initial position: Lat=%.6f, Lon=%.6f\n", p.lat, p.lon)
	log.Println("Press Ctrl+C to stop...")

	for {
		select {
		case <-ticker.C:
			p.publishLocation(vehicleId)
		}
	}
}

// PublishOnce publishes a single location update (for testing)
func (p *MockPublisher) PublishOnce(vehicleId string) error {
	return p.publishLocation(vehicleId)
}

func (p *MockPublisher) publishLocation(vehicleId string) error {
	// Random movement: ±0.0001 degrees (~11 meters)
	// This simulates vehicle moving around the geofence area
	p.lat += (rand.Float64() - 0.5) * 0.0002
	p.lon += (rand.Float64() - 0.5) * 0.0002

	payload := model.MQTTLocationStruct{
		VehicleId: vehicleId,
		Latitude:  p.lat,
		Longitude: p.lon,
		Timestamp: time.Now().Unix(),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleId)
	token := p.client.Publish(topic, 0, false, data)

	token.Wait()

	if token.Error() != nil {
		log.Fatalf("Publish failed: %v\n", token.Error())
		return token.Error()
	}

	log.Printf("Published: %s @ %.6f, %.6f (timestamp: %d)\n",
		vehicleId, p.lat, p.lon, payload.Timestamp)

	return nil
}

// SetPosition manually sets the vehicle position (for testing geofence)
func (p *MockPublisher) SetPosition(lat, lon float64) {
	p.lat = lat
	p.lon = lon

	log.Println("Position set to: Lat=%.6f, Lon=%.6f\n", lat, lon)
}
