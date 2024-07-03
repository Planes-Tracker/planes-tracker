package datasource

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/HewlettPackard/structex"
	"github.com/LockBlock-dev/planes-tracker/internal/client"
	"github.com/LockBlock-dev/planes-tracker/internal/types"
	"github.com/LockBlock-dev/planes-tracker/internal/types/adsbexchange"
	"github.com/klauspost/compress/zstd"
	"google.golang.org/protobuf/proto"
)

func apiBufferToString(slice []byte) string {
	n := bytes.IndexByte(slice[:], 0)

	if n != -1 {
		return string(slice[:n])
	}

	return string(slice[:])
}

type ADSBExchangeDataSource struct {
	Client  *client.ADSBExchangeClient
	Decoder *zstd.Decoder
}

func NewADSBExchangeDataSource() (*ADSBExchangeDataSource, error) {
	client, err := client.NewADSBExchangeClient()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize HTTP client: %w",
			err,
		)
	}

	decoder, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize zstd decoder: %w",
			err,
		)
	}
	
	return &ADSBExchangeDataSource{
		Client:  client,
		Decoder: decoder,
	}, nil
}

func (src *ADSBExchangeDataSource) Name() string {
	return "ADS-B Exchange HTTP"
}

func (src *ADSBExchangeDataSource) FetchFlights(ch chan<- types.FlightRecord, location *types.Coordinates, radius *types.Radius) (int, error) {
	resp, err := src.Client.FetchFlights(location, radius)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected HTTP status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	decompressed, err := src.Decoder.DecodeAll(body, nil)
	if err != nil {
		return 0, err
	}

	var data adsbexchange.ADSBExchangeData

	reader := bytes.NewReader(decompressed)

	if err := binary.Read(reader, binary.LittleEndian, &data); err != nil {
		return 0, err
	}

	if _, err := reader.Seek(int64(data.ElementSize), io.SeekStart); err != nil {
		return 0, err
	}

	for reader.Len() >= int(data.ElementSize) {
		var a adsbexchange.BinaryAircraft

		if reader.Len() < int(data.ElementSize) {
			return 0, fmt.Errorf("api response buffer does not contain enough data")
		}

		if err := structex.Decode(reader, &a); err != nil {
			return 0, err
		}

		// https://github.com/ADSBexchange/tar1090/blob/master/html/formatter.js#L402
		flight := types.Flight{
			Registration: proto.String(strings.TrimRight(apiBufferToString(a.Registration[:]), " ")),
			Callsign: proto.String(strings.TrimRight(apiBufferToString(a.Callsign[:]), " ")),
			Latitude: proto.Float32(float32(a.Lat) / 1e6),
			Longitude: proto.Float32(float32(a.Lon) / 1e6),
			Altitude: proto.Int32(int32(a.BaroAlt) * 25),
			Track: proto.Int32(int32(a.Track) / 90),
			Speed: proto.Int32(int32(a.GS) / 10),
			VerticalSpeed: proto.Int32(int32(a.BaroRate) * 8),
			OnGround: proto.Bool(a.BaroAlt <= 0),
			SquawkCode: proto.String(fmt.Sprintf("%04X", a.Squawk)),
			Model: proto.String(apiBufferToString(a.TypeCode[:])),
			ICAOAddress: proto.String(fmt.Sprintf("%06X", a.Hex)),
		}

		ch <- types.FlightRecord{
			Datasource: src.Name(),
			Flight: flight,
		}
	}

	return int(data.ResultCount), nil
}
