package types

import (
	"math"
)

type Coordinates struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type Radius struct {
	// kilometers
	Distance uint32 `json:"distance"`
}

func (r *Radius) AsMeters() uint32 {
	return r.Distance * 1000
}

func (r *Radius) AsKiloMeters() uint32 {
	return r.Distance
}

func (r *Radius) AsNauticalMiles() uint32 {
	return r.AsMeters() / 1852
}

type BoundingBox struct {
	NorthWest Coordinates
	SouthEast Coordinates
}

func NewBoundingBox(lat float64, lon float64, radiusDistance float64) *BoundingBox {
	// IUGG mean radius in km
	const earthRadius = 6371.0

	// Convert radiusDistance to radians
	dLat := radiusDistance / earthRadius
	dLon := radiusDistance / (earthRadius * math.Cos(math.Pi*lat/180.0))

	box := &BoundingBox{
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
