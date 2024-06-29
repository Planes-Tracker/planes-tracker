package app

import (
	"fmt"
	"log"
	"time"

	"github.com/LockBlock-dev/planes-tracker/internal/database"
	"github.com/LockBlock-dev/planes-tracker/internal/datasource"
	"github.com/LockBlock-dev/planes-tracker/internal/entities"
	"github.com/LockBlock-dev/planes-tracker/internal/types"
	"github.com/mattn/go-sqlite3"
)

type App struct {
	DB          *database.Database
	Config      *types.Config
	Ticker      *time.Ticker
	DataSources []types.DataSource
	Channel     chan types.FlightRecord
}

func NewApp() (*App, error) {
	config, err := types.NewConfigFromFile("./config.json")
	if err != nil {
		return nil, fmt.Errorf(
			"failed to load configuration file: %w",
			err,
		)
	}
	
	fr24DataSource, err := datasource.NewFR24DataSource()
	if err != nil {
		log.Println(err)
	}

	adsbexchangeDataSource, err := datasource.NewADSBExchangeDataSource()
	if err != nil {
		log.Println(err)
	}

	app := App{
		Config: config,
		Ticker: time.NewTicker(time.Duration(config.PollRate) * time.Second),
		DataSources: []types.DataSource{
			fr24DataSource,
			adsbexchangeDataSource,
		},
		Channel: make(chan types.FlightRecord),
	}

	verbose := false
    db, err := database.InitDB(app.Config.DatabaseName, verbose)
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

func (app *App) watch() {
	for flightRecord := range app.Channel {
		var flight entities.Flight

		result := app.DB.Where(
			&entities.Flight{
				Registration: flightRecord.Flight.Registration,
				Callsign:     flightRecord.Flight.Callsign,
				ICAOAddress:  flightRecord.Flight.ICAOAddress,
			},
		).Attrs(
			&entities.Flight{
				Flight:      flightRecord.Flight.Flight,
				Origin:      flightRecord.Flight.Origin,
				Destination: flightRecord.Flight.Destination,
				DivertedTo:  flightRecord.Flight.DivertedTo,
				Model:       flightRecord.Flight.Model,
			},
		).FirstOrCreate(
			&flight,
		)

		if result.Error != nil {		
			sqlErr, ok := result.Error.(sqlite3.Error)
			
			if (ok && sqlErr.Code != sqlite3.ErrConstraint) || !ok {
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

		scheduleForUpdateIfNeeded("flight", flight.Flight, flightRecord.Flight.Flight)
		scheduleForUpdateIfNeeded("origin", flight.Origin, flightRecord.Flight.Origin)
		scheduleForUpdateIfNeeded("destination", flight.Destination, flightRecord.Flight.Destination)
		scheduleForUpdateIfNeeded("diverted_to", flight.DivertedTo, flightRecord.Flight.DivertedTo)

		if len(updates) > 0 {
			app.DB.Model(&flight).Updates(updates)
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
			_, ok := err.(sqlite3.Error)

			if !ok {
				log.Printf("Failed to append flight point to flight with id %d, error: %v\n", flight.FlightId, err)
			}

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
