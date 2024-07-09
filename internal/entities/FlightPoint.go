package entities

import "time"

type FlightPoint struct {
	FlightId      uint
	FlightPointId uint `gorm:"primarykey"`

	CreatedAt     time.Time
	UpdatedAt     time.Time
	Latitude      *float32 `gorm:"column:latitude"`
	Longitude     *float32 `gorm:"column:longitude"`
	Altitude      *int32   `gorm:"column:altitude"`
	Track         *int32   `gorm:"column:track"`
	Speed         *int32   `gorm:"column:speed;default:0"`
	VerticalSpeed *int32   `gorm:"column:vertical_speed"`
	OnGround      *bool    `gorm:"column:on_ground;default:false"`
	SquawkCode    *string  `gorm:"column:squawk"`
}
