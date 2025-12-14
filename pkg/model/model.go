package models

import (
	"time"
)

type MQTTLocationStruct struct {
	VehicleId string  `json:"vehicle_id" gorm:"column:vehicle_id"`
	Latitude  float64 `json:"latitude" gorm:"column:latitude"`
	Longitude float64 `json:"longitude" gorm:"column:longitude"`
	Timestamp int64   `json:"timestamp" gorm:"column:timestamp"`
}

type VehicleLocation struct {
	Id int64 `json:"id" gorm:"column:id;primaryKey"`
	MQTTLocationStruct
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (MQTTLocationStruct) TableName() string {
	return "vehicle_locations"
}

type BusStation struct {
	Id        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Latitude  float64   `gorm:"column:latitude"`
	Longitude float64   `gorm:"column:longitude"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (BusStation) TableName() string {
	return "bus_stations"
}
