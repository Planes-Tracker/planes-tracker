package entities

import "time"

type Flight struct {
	FlightId uint `gorm:"primarykey"`

	CreatedAt    time.Time
	UpdatedAt    time.Time
	Registration *string `gorm:"column:registration;uniqueIndex:idx_unique_flight"`
	Flight       *string `gorm:"column:flight"`
	Callsign     *string `gorm:"column:callsign;uniqueIndex:idx_unique_flight"`
	Origin       *string `gorm:"column:origin"`
	Destination  *string `gorm:"column:destination"`
	DivertedTo   *string `gorm:"column:diverted_to"`
	Model        *string `gorm:"column:model"`
	ICAOAddress  *string `gorm:"column:icao_address;type:char(6);uniqueIndex:idx_unique_flight"`

	FlightPoints []FlightPoint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
