package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/LockBlock-dev/planes-tracker/internal/fr24"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Coordinates struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type Config struct {
	Location		Coordinates	`json:"location"`
	RadiusDistance	int			`json:"radiusDistance"`
	PollRate		int			`json:"pollRate"`
	Debug			bool		`json:"debug"`
}


type BoundingBox struct {
	NorthWest Coordinates
	SouthEast Coordinates
}

type Plane struct {
    ID           uint32
    Registration string
    Flight       *string
    Callsign     *string
    Origin       *string
    Destination  *string
    Latitude     float32
    Longitude    float32
    Altitude     int32
    Bearing      int32
    Speed        int32
    RateOfClimb  int32
    IsOnGround   bool
    SqawkCode    *string
    Model        *string
    ModeSCode    *string
    EnteredAt    time.Time
    LeftAt       time.Time
}

var config Config

func NewBoundingBox(lat float64, lon float64, radiusDistance float64) BoundingBox {
	// IUGG mean radius in km
	const earthRadius = 6371.0

	// Convert radiusDistance to radians
	dLat := radiusDistance / earthRadius
	dLon := radiusDistance / (earthRadius * math.Cos(math.Pi*lat/180.0))

	box := BoundingBox{
		NorthWest: Coordinates{
			Lat: lat + (dLat * 180.0 / math.Pi),
			Lon: lon - (dLon * 180.0 / math.Pi),
		},
		SouthEast: Coordinates{
			Lat: lat - (dLat * 180.0 / math.Pi),
			Lon: lon + (dLon * 180.0 / math.Pi),
		},
	}

	return box
}

func fetchFlights(box BoundingBox) (*fr24.LiveFeedResponse, error) {
	req := &fr24.LiveFeedRequest{
		Bounds: &fr24.LocationBoundaries{
			North: float32(box.NorthWest.Lat),
			South: float32(box.SouthEast.Lat),
			West: float32(box.NorthWest.Lon),
			East: float32(box.SouthEast.Lon),
		},
		Settings: &fr24.VisibilitySettings{
			Sources: []fr24.DataSource{
				fr24.DataSource_ADSB,
				fr24.DataSource_MLAT,
				fr24.DataSource_FLARM,
				fr24.DataSource_FAA,
				fr24.DataSource_ESTIMATED,
				fr24.DataSource_SATELLITE,
				fr24.DataSource_OTHER_DATA_SOURCE,
				fr24.DataSource_UAT,
				fr24.DataSource_SPIDERTRACKS,
				fr24.DataSource_AUS,
			},
			Services: []fr24.Service{
				fr24.Service_PASSENGER,
				fr24.Service_CARGO,
				fr24.Service_MILITARY_AND_GOVERNMENT,
				fr24.Service_BUSINESS_JETS,
				fr24.Service_GENERAL_AVIATION,
				fr24.Service_HELICOPTERS,
				fr24.Service_LIGHTER_THAN_AIR,
				fr24.Service_GLIDERS,
				fr24.Service_DRONES,
				fr24.Service_GROUND_VEHICLES,
				fr24.Service_OTHER_SERVICE,
				fr24.Service_NON_CATEGORIZED,
			},
			TrafficType: fr24.TrafficType_ALL.Enum(),
			OnlyRestricted: proto.Bool(false),
		},
		HighlightMode: proto.Bool(false),
		Stats: proto.Bool(false),
		Limit: proto.Uint32(10),
		Maxage: proto.Uint32(14400),
		FieldMask: &fieldmaskpb.FieldMask{
			Paths: []string{"flight", "reg", "route", "type", "schedule", "icao_address"},
		},
		RestrictionMode: fr24.RestrictionVisibility_NOT_VISIBLE.Enum(),
	}

	creds := credentials.NewTLS(&tls.Config{});

	conn, err := grpc.Dial("data-feed.flightradar24.com:443", grpc.WithTransportCredentials(creds))
	if err != nil {
		conn.Close()
		return nil, err
	}
	defer conn.Close()

	feed := fr24.NewFeedClient(conn)

	headers := metadata.New(map[string]string{
		"fr24-device-id": "web",
	})

	ctx := metadata.NewOutgoingContext(context.Background(), headers)

	resp, err := feed.LiveFeed(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil;
}

func fetchAndSave(box BoundingBox, db *gorm.DB) {
	resp, err := fetchFlights(box)
	if err != nil {
		log.Fatal(err)
	}

	if config.Debug {
		log.Println(fmt.Sprintf("Total planes currently in range: %d", len(resp.Flights)))
	}
	
	for _, flight := range resp.Flights {
		var speed int32 = 0
		var vSpeed int32 = 0
		var squawk string
		var sqawkPtr *string
		var icaoAddress string
		var icaoAddressPtr *string

		if flight.Speed != nil {
			speed = *flight.Speed
		}

		if flight.ExtraInfo.Vspeed != nil {
			vSpeed = *flight.ExtraInfo.Vspeed
		}

		if flight.ExtraInfo.OperatedById != nil {
			icaoAddress = fmt.Sprintf("%X", *flight.ExtraInfo.OperatedById)
			icaoAddressPtr = &icaoAddress
		}

		if flight.ExtraInfo.Squawk != nil {
			squawk = string(*flight.ExtraInfo.Squawk)
			sqawkPtr = &squawk
		}

		var plane Plane

		if err := db.Where("id = ?", *flight.FlightId).First(&plane).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&Plane{
					ID: *flight.FlightId,
					Registration: *flight.ExtraInfo.Reg,
					Flight: flight.ExtraInfo.Flight,
					Callsign: flight.Callsign,
					Origin: flight.ExtraInfo.Route.FromIata,
					Destination: flight.ExtraInfo.Route.ToIata,
					Latitude: *flight.Lat,
					Longitude: *flight.Lon,
					Altitude: *flight.Alt,
					Bearing: *flight.Track,
					Speed: speed,
					RateOfClimb: vSpeed,
					IsOnGround: *flight.OnGround,
					SqawkCode: sqawkPtr,
					Model: flight.ExtraInfo.Type,
					ModeSCode: icaoAddressPtr,
					EnteredAt: time.Unix(int64(*flight.Timestamp), 0),
					LeftAt: time.Unix(int64(*flight.Timestamp), 0),
				})
				continue
			} else {
				log.Fatal(err)
			}
		}

		plane.LeftAt = time.Unix(int64(*flight.Timestamp), 0)
		db.Save(&plane)
	}
}

func initFetching(box BoundingBox, db *gorm.DB, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Duration(config.PollRate) * time.Second)
	defer ticker.Stop()

	if config.Debug {
		log.Println("Fetching for planes...")
	}

	fetchAndSave(box, db)

	for range ticker.C {
		if config.Debug {
			log.Println("Fetching for planes...")
		}

		fetchAndSave(box, db)
	}

	defer wg.Done()
}

func main()  {
	configFile, err := os.Open("./config.json")
    if err != nil {
        configFile, err = os.Open("../../config.json")
		if err != nil {
			log.Fatal(fmt.Errorf("Cannot find/open config.json file: %w", err))
		}
    }
    defer configFile.Close()

    decoder := json.NewDecoder(configFile)
    config = Config{}
    err = decoder.Decode(&config)
    if err != nil {
        log.Fatal(fmt.Errorf("Cannot parse config.json file: %w", err))
    }
	
	db, err := gorm.Open(sqlite.Open("planes.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to connect database: %w", err))
	}

	err = db.AutoMigrate(&Plane{})
	if err != nil {
		log.Fatal(err)
	}
  
	box := NewBoundingBox(config.Location.Lat, config.Location.Lon, float64(config.RadiusDistance))

	var wg sync.WaitGroup
	wg.Add(1)

	go initFetching(box, db, &wg)

	println("Planes tracker started!")

	wg.Wait()
}
