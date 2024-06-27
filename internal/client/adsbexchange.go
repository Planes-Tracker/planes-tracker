package client

import (
	"fmt"
	"net/http"

	"github.com/LockBlock-dev/planes-tracker/internal/types"
)

type ADSBExchangeClient struct {
	client  *http.Client
}

func NewADSBExchangeClient() (*ADSBExchangeClient, error) {
	return &ADSBExchangeClient{
		client:  &http.Client{},
	}, nil
}

func (c *ADSBExchangeClient) FetchFlights(location *types.Coordinates, radius *types.Radius) (*http.Response, error) {
	url := fmt.Sprintf(
		"https://globe.adsbexchange.com/re-api/?binCraft&zstd&circle=%f,%f,%d",
		location.Lat,
		location.Lon,
		// ADS-B Exchange API takes Nautical Miles
		radius.AsNauticalMiles(),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:113.0) Gecko/20100101 Firefox/113.0")
	req.Header.Set("Referer", "https://globe.adsbexchange.com/")

	return c.client.Do(req)
}
