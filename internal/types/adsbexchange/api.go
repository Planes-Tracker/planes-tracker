package adsbexchange

// https://github.com/ADSBexchange/readsb/blob/eea98b13b453779b9fe75b6783f4e6b129f2dff0/api.c#L615-L653
type ADSBExchangeData struct {
	Timestamp                 int64
	ElementSize               uint32
	AircraftWithPositionCount uint32
	Index                     uint32
	South                     int16
	West                      int16
	North                     int16
	East                      int16
	MessageCount              uint32
	ResultCount               uint32
	Dummy                     int32
	BinCraftVersion           uint32
	MessageRate               uint32
}
