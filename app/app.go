package app

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/LockBlock-dev/planes-tracker/config"
	"github.com/LockBlock-dev/planes-tracker/internal/database"
	"github.com/LockBlock-dev/planes-tracker/internal/datasource"
	"github.com/LockBlock-dev/planes-tracker/internal/entities"
	"github.com/LockBlock-dev/planes-tracker/internal/types"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type App struct {
	DB          *database.Database
	Config      *types.Config
	Ticker      *time.Ticker
	DataSources []types.DataSource
	Channel     chan types.FlightRecord
}

func NewApp() (*App, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print(fmt.Errorf("failed to load .env file: %w", err))
		log.Println("Assuming we are running inside Docker...")
	}

	appConfig := config.MakeConfigFromEnv()

	fr24DataSource, err := datasource.NewFR24DataSource()
	if err != nil {
		log.Println(err)
	}

	adsbexchangeDataSource, err := datasource.NewADSBExchangeDataSource()
	if err != nil {
		log.Println(err)
	}

	app := App{
		Config: appConfig,
		Ticker: time.NewTicker(time.Duration(appConfig.PollRate) * time.Second),
		DataSources: []types.DataSource{
			fr24DataSource,
			adsbexchangeDataSource,
		},
		Channel: make(chan types.FlightRecord),
	}

	db, err := database.InitDB(app.Config.Debug)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize the database: %w",
			err,
		)
	}

	app.DB = db

	return &app, nil
}

func (app *App) Start() {
	if app.Config.Debug {
		log.Println("Tracker started!")
	}

	go app.watch()

	app.run()

	for range app.Ticker.C {
		app.run()
	}
}

func (app *App) Stop() {
	if err := app.DB.Close(); err != nil {
		log.Printf("Failed to close the database: %v\n", err)
	}

	close(app.Channel)
}

func preventEmptyString(s *string) *string {
	if s != nil && *s == "" {
		return nil
	}

	return s
}

func (app *App) watch() {
	for flightRecord := range app.Channel {
		var flight entities.Flight

		result := app.DB.Where(
			&entities.Flight{
				Registration: preventEmptyString(flightRecord.Flight.Registration),
				Callsign:     preventEmptyString(flightRecord.Flight.Callsign),
				ICAOAddress:  preventEmptyString(flightRecord.Flight.ICAOAddress),
			},
		).Attrs(
			&entities.Flight{
				Flight:      preventEmptyString(flightRecord.Flight.Flight),
				Origin:      preventEmptyString(flightRecord.Flight.Origin),
				Destination: preventEmptyString(flightRecord.Flight.Destination),
				DivertedTo:  preventEmptyString(flightRecord.Flight.DivertedTo),
				Model:       preventEmptyString(flightRecord.Flight.Model),
			},
		).FirstOrCreate(
			&flight,
		)

		if result.Error != nil {
			if result.Error != nil && errors.Is(result.Error, gorm.ErrDuplicatedKey) {
				log.Printf("Failed to insert flights into the database: %v\n", result.Error)
			}

			continue
		}

		updates := map[string]interface{}{}

		scheduleForUpdateIfNeeded := func(fieldName string, currentValue interface{}, newValue interface{}) {
			if currentValue == nil && newValue != currentValue {
				updates[fieldName] = newValue
			}
		}

		scheduleForUpdateIfNeeded("flight", flight.Flight, preventEmptyString(flightRecord.Flight.Flight))
		scheduleForUpdateIfNeeded("origin", flight.Origin, preventEmptyString(flightRecord.Flight.Origin))
		scheduleForUpdateIfNeeded("destination", flight.Destination, preventEmptyString(flightRecord.Flight.Destination))
		scheduleForUpdateIfNeeded("diverted_to", flight.DivertedTo, preventEmptyString(flightRecord.Flight.DivertedTo))

		if len(updates) > 0 {
			app.DB.Model(&flight).Updates(updates)
		}

		var recentFlightPoint entities.FlightPoint

		app.DB.Where(
			&entities.FlightPoint{FlightId: flight.FlightId},
		).Order(
			clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true},
		).First(&recentFlightPoint)

		if recentFlightPoint.FlightPointId != 0 {
			timeElapsed := time.Since(recentFlightPoint.CreatedAt)
			if timeElapsed < time.Duration(app.Config.PollRate/2)*time.Second {
				// skip saving duplicated flight point
				continue
			}
		}

		err := app.DB.Model(&flight).Association("FlightPoints").Append(
			&entities.FlightPoint{
				Latitude:      flightRecord.Flight.Latitude,
				Longitude:     flightRecord.Flight.Longitude,
				Altitude:      flightRecord.Flight.Altitude,
				Track:         flightRecord.Flight.Track,
				Speed:         flightRecord.Flight.Speed,
				VerticalSpeed: flightRecord.Flight.VerticalSpeed,
				OnGround:      flightRecord.Flight.OnGround,
				SquawkCode:    flightRecord.Flight.SquawkCode,
			},
		)

		if err != nil {
			log.Printf("Failed to append flight point to flight with id %d, error: %v\n", flight.FlightId, err)

			continue
		}
	}
}

func (app *App) run() {
	for _, dataSource := range app.DataSources {
		if app.Config.Debug {
			log.Printf("Fetching flights from %s...\n", dataSource.Name())
		}

		count, err := dataSource.FetchFlights(app.Channel, &app.Config.Location, &app.Config.Radius)
		if err != nil {
			log.Printf("Error fetching flights: %v\n", err)
			// Continue with the next data source
			continue
		}

		if app.Config.Debug {
			log.Printf("Got %d flights from %s\n", count, dataSource.Name())
		}
	}
}
