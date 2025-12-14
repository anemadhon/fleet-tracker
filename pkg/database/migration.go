package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(migrationsPath string) error {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("could not resolve absolute path: %w", err)
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory not found: %s", absPath)
	}

	log.Printf("Migrations path: %s", absPath)

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("could not get sql.DB from gorm: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")

	return nil
}

func RollbackMigration(migrationsPath string, steps int) error {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("could not resolve absolute path: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("could not get sql.DB from gorm: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	sourceURL := fmt.Sprintf("file://%s", filepath.ToSlash(absPath))
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	if err := m.Steps(-steps); err != nil {
		return fmt.Errorf("could not rollback migrations: %w", err)
	}

	log.Println("Rolled back %d migration(s)", steps)
	return nil
}

func ForceMigrationVersion(migrationsPath string, version uint) error {
    absPath, err := filepath.Abs(migrationsPath)
    if err != nil {
        return fmt.Errorf("could not resolve absolute path: %w", err)
    }

    sqlDB, err := DB.DB()
    if err != nil {
        return fmt.Errorf("could not get sql.DB from gorm: %w", err)
    }

    driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("could not create postgres driver: %w", err)
    }

    sourceURL := fmt.Sprintf("file://%s", filepath.ToSlash(absPath))
    m, err := migrate.NewWithDatabaseInstance(
        sourceURL,
        "postgres",
        driver,
    )
    if err != nil {
        return fmt.Errorf("could not create migrate instance: %w", err)
    }
    if err := m.Force(int(version)); err != nil {
        return fmt.Errorf("could not force migration version: %w", err)
    }

    log.Printf("Forced migration version to %d (dirty reset)", version)
    return nil
}

