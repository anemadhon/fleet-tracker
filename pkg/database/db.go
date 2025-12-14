package database

import (
	"log"
	"os"
	"time"
	"tj/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	gormLogger := logger.New(
		log.New(os.Stdout, "GORM ", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(postgres.Open(config.Cfg.PostgresDSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}

	sqlDB, _ := db.DB()

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	log.Println("PostgreSQL connected")

	return nil
}
