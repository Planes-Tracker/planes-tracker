package types

type Config struct {
	PollRate int  `json:"pollRate"`
	Debug    bool `json:"debug"`

	Location Coordinates `json:"location"`
	Radius   Radius      `json:"radius"`
}
