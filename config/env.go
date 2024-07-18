package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/LockBlock-dev/planes-tracker/internal/types"
)

func GetEnvWithDefault(key string, defaultVal string) string {
	val := os.Getenv(key)

	if val == "" {
		return defaultVal
	}

	return val
}

func MakeConfigFromEnv() *types.Config {
	// Maybe use Viper in the future? So we can use getters like GetBool
	pollRate, _ := strconv.Atoi(GetEnvWithDefault("TRACKER_POLL_RATE", "10"))
	isDebug, _ := strconv.ParseBool(GetEnvWithDefault("TRACKER_DEBUG", "false"))
	latitude, _ := strconv.ParseFloat(os.Getenv("TRACKER_LOCATION_LATITUDE"), 64)
	longitude, _ := strconv.ParseFloat(os.Getenv("TRACKER_LOCATION_LONGITUDE"), 64)
	radiusDistance, _ := strconv.Atoi(os.Getenv("TRACKER_RADIUS_DISTANCE"))

	return &types.Config{
		PollRate: pollRate,
		Debug:    isDebug,
		Location: types.Coordinates{
			Lat: latitude,
			Lon: longitude,
		},
		Radius: types.Radius{
			Distance: uint32(radiusDistance),
		},
	}
}

func DatabaseDSN() string {
	postgresUser := GetEnvWithDefault("POSTGRES_USER", "postgres")

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		GetEnvWithDefault("POSTGRES_HOST", "postgres"),
		postgresUser,
		os.Getenv("POSTGRES_PASSWORD"),
		GetEnvWithDefault("POSTGRES_DB", postgresUser),
		GetEnvWithDefault("POSTGRES_PORT", "5432"),
		GetEnvWithDefault("POSTGRES_TIMEZONE", "UTC"),
	)
}
