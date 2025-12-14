package main

import (
	"flag"
	"log"
	"tj/config"

	db "tj/pkg/database"
)

func main() {
	action := flag.String("action", "up", "Migration action: up, down, rollback, force")
	steps := flag.Int("steps", 1, "Number of steps to rollback")
	version := flag.Int("version", 0, "Version to force")

	flag.Parse()

	config.Load()

	if err := db.Connect(); err != nil {
		log.Fatalf("Postgres init error: %v", err)
	}

	switch *action {
	case "up":
		if err := db.RunMigrations("../../../migrations"); err != nil {
			log.Fatal(err)
		}

	case "down", "rollback":
		if err := db.RollbackMigration("../../../migrations", *steps); err != nil {
			log.Fatal(err)
		}

	case "force":
		if *version == 0 {
			log.Fatal("force action requires -version")
		}
		if err := db.ForceMigrationVersion("../../../migrations", uint(*version)); err != nil {
			log.Fatal(err)
		}
	}
}
