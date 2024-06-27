package client

import (
	"context"
	"crypto/tls"

	"github.com/LockBlock-dev/planes-tracker/internal/types"
	"github.com/LockBlock-dev/planes-tracker/internal/types/fr24"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type FR24Client struct {
	conn    *grpc.ClientConn
	client  fr24.FeedClient
	context context.Context
}

func NewFR24Client() (*FR24Client, error) {
	creds := credentials.NewTLS(&tls.Config{})

	conn, err := grpc.Dial(
		"data-feed.flightradar24.com:443",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return nil, err
	}

	client := fr24.NewFeedClient(conn)
	headers := metadata.New(map[string]string{
		"fr24-device-id": "web",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), headers)

	return &FR24Client{
		conn:    conn,
		client:  client,
		context: ctx,
	}, nil
}

func (c *FR24Client) FetchFlights(location *types.Coordinates, radius *types.Radius) (*fr24.NearestFlightsResponse, error) {
	req := &fr24.NearestFlightsRequest{
		Location: &fr24.Geolocation{
			Lat: float32(location.Lat),
			Lon: float32(location.Lon),
		},
		Radius: proto.Uint32(radius.AsMeters()),
		Limit: proto.Uint32(500),
	}

	return c.client.NearestFlights(c.context, req)
}

func (c *FR24Client) FetchFlight(flightId uint32) (*fr24.LiveFeedResponse, error) {
	req := &fr24.LiveFeedRequest{
		SelectedFlightIds: []uint32{flightId},
		Limit: proto.Uint32(1),
	}

	return c.client.LiveFeed(c.context, req)
}

func (c *FR24Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
