package entities

import "time"

type Flight struct {
	FlightId uint `gorm:"primarykey"`

	CreatedAt    time.Time
	UpdatedAt    time.Time
	Registration *string `gorm:"column:registration;uniqueIndex:idx_unique_flight"`
	Flight       *string `gorm:"column:flight;uniqueIndex:idx_unique_flight"`
	Callsign     *string `gorm:"column:callsign;uniqueIndex:idx_unique_flight"`
	Origin       *string `gorm:"column:origin;uniqueIndex:idx_unique_flight"`
	Destination  *string `gorm:"column:destination;uniqueIndex:idx_unique_flight"`
	DivertedTo   *string `gorm:"column:diverted_to"`
	Model        *string `gorm:"column:model"`
	ICAOAddress  *string `gorm:"column:icao_address;type:char(6);uniqueIndex:idx_unique_flight"`

	FlightPoints []FlightPoint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
