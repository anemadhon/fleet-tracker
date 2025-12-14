package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"tj/config"
	db "tj/pkg/database"
	handler "tj/services/api/internal/controller"
)

func main() {
	config.Load()

	if err := db.Connect(); err != nil {
		log.Fatalf("Postgres init error: %v", err)
	}

	r := gin.Default()
	vh := handler.NewVehicleHandler(db.DB)

	r.GET("/vehicles/:vehicle_id/location", vh.GetLastLocation)
	r.GET("/vehicles/:vehicle_id/history", vh.GetHistory)

	if err := r.Run(":8093"); err != nil {
		log.Fatalf("Gin run error: %v", err)
	}
}
