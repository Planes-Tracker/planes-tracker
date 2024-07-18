package database

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/LockBlock-dev/planes-tracker/config"
	"github.com/LockBlock-dev/planes-tracker/internal/entities"
)

type Database struct {
	*gorm.DB
}

func InitDB(verbose bool) (*Database, error) {
	var dbLogger logger.Interface

	if verbose {
		dbLogger = logger.Default
	} else {
		dbLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(
		postgres.Open(
			config.DatabaseDSN(),
		),
		&gorm.Config{
			Logger:         dbLogger,
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
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

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
