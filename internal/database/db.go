package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/LockBlock-dev/planes-tracker/internal/entities"
)

type Database struct {
	*gorm.DB
}

func getEnvWithDefault(key string, defaultVal string) string {
	val := os.Getenv(key)

	if val == "" {
		return defaultVal
	}

	return val
}

func InitDB(verbose bool) (*Database, error) {
	var dbLogger logger.Interface
	
	if verbose {
		dbLogger = logger.Default
	} else {
		dbLogger = logger.Default.LogMode(logger.Silent)
	}

	postgresUser := getEnvWithDefault("POSTGRES_USER", "postgres")
	postgresDatabase := getEnvWithDefault("POSTGRES_DB", postgresUser)

	db, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf(
				"host=postgres user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=%s",
				postgresUser,
				os.Getenv("POSTGRES_PASSWORD"),
				postgresDatabase,
				getEnvWithDefault("POSTGRES_TIMEZONE", "UTC"),
			),
		),
		&gorm.Config{
			Logger: dbLogger,
			TranslateError: true,
		},
	)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(entities.Flight{}); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(entities.FlightPoint{}); err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxIdleTime(time.Hour)

	return &Database{
		db,
	}, err
}

func (database *Database) Close() error {
	sqlDB, err := database.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
