package datasource

import (
	"fmt"

	"github.com/LockBlock-dev/planes-tracker/internal/client"
	"github.com/LockBlock-dev/planes-tracker/internal/types"
	"google.golang.org/protobuf/proto"
)

type FR24DataSource struct {
	Client *client.FR24Client
	name   string
}

func NewFR24DataSource() (*FR24DataSource, error) {
	client, err := client.NewFR24Client()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize gRPC client: %w",
			err,
		)
	}

	return &FR24DataSource{
		Client: client,
		name:   "Flightradar24 gRPC",
	}, nil
}

func (src *FR24DataSource) Name() string {
	return src.name
}

func (src *FR24DataSource) FetchFlights(ch chan<- types.FlightRecord, location *types.Coordinates, radius *types.Radius) (int, error) {
	resp, err := src.Client.FetchFlights(location, radius)
	if err != nil {
		return 0, err
	}

	// fmt.Println(prototext.Format(resp))

	for _, nearbyFlight := range resp.Flights {
		data := nearbyFlight.Flight
		flight := types.Flight{
			Callsign:  data.Callsign,
			Latitude:  data.Lat,
			Longitude: data.Lon,
			Altitude:  data.Alt,
			Track:     data.Track,
			Speed:     data.Speed,
			OnGround:  data.OnGround,
		}
		extraInfo := data.ExtraInfo

		if extraInfo != nil {
			flight.Registration = extraInfo.Reg
			flight.Flight = extraInfo.Flight
			flight.VerticalSpeed = extraInfo.Vspeed
			flight.Model = extraInfo.Type

			if extraInfo.Route != nil {
				flight.Origin = extraInfo.Route.From
				flight.Destination = extraInfo.Route.To
				flight.DivertedTo = extraInfo.Route.DivertedTo
			}
		}

		if data.FlightId != nil {
			resp, err := src.Client.FetchFlight(*data.FlightId)
			if err == nil && len(resp.SelectedFlights) > 0 {
				extraInfo := resp.SelectedFlights[0].ExtraInfo

				if extraInfo != nil {
					if extraInfo.IcaoAddress != nil {
						flight.ICAOAddress = proto.String(fmt.Sprintf("%06X", *extraInfo.IcaoAddress))
					}

					if extraInfo.Squawk != nil {
						flight.SquawkCode = proto.String(fmt.Sprintf("%04X", *extraInfo.Squawk))
					}
				}
			}
		}

		ch <- types.FlightRecord{
			Datasource: src.Name(),
			Flight:     flight,
		}
	}

	return len(resp.Flights), nil
}
