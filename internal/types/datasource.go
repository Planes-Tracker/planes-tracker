package types

type DataSource interface {
	Name() string
	FetchFlights(ch chan<- FlightRecord, location *Coordinates, radius *Radius) (int, error)
}
