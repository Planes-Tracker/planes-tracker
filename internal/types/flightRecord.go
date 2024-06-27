package types

type FlightRecord struct {
    Datasource string
    Flight
}

type Flight struct {
    Registration  *string
    Flight        *string
    Callsign      *string
    Origin        *string
    Destination   *string
    DivertedTo    *string
    Latitude      *float32
    Longitude     *float32
    Altitude      *int32
    Track         *int32
    Speed         *int32
    VerticalSpeed *int32
    OnGround      *bool
    SquawkCode    *string
    Model         *string
    ICAOAddress   *string
}
