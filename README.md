# Planes tracker

Saves all the planes flying above the designated area in a database.

See the [changelog](/CHANGELOG.md) for the latest updates.

## Table of content

-   [**Installation**](#installation)
-   [**Compiling from source**](#compiling-from-source)
-   [**Configuration**](#configuration)
-   [**Config details**](#config-details)
-   [**Usage**](#usage)
-   [**Credits**](#credits)
-   [**Copyright**](#copyright)

## Installation

You can [use Docker](#usage) or install this app manually. Here's how:

-   Download [go](https://go.dev/dl/) (go 1.20 required).
-   Download or clone the project.
-   Download the binary from the [Releases](../../releases) or [build it](#compiling-from-source) yourself.
-   [Configure the app](#configuration).

## Compiling from source

You can either:

-   Use the [`Makefile`](/Makefile) by running `make` in the project directory.
-   Use `go build` in the project directory.

## Configuration

The configuration details must be set inside a `.env` file at the root of the project. An exemple is provided inside [`.env.example`](/.env.example).

## Config details

| Item                         | Values                  | Meaning                                                           |
| ---------------------------- | ----------------------- | ----------------------------------------------------------------- |
| `TRACKER_POLL_RATE`          | `number`                | Tracker poll rate (in seconds)                                    |
| `TRACKER_DEBUG`              | `boolean`               | Enable debug logs                                                 |
| `TRACKER_LOCATION_LATITUDE`  | `floating point number` | Latitude of the center of the area to be covered                  |
| `TRACKER_LOCATION_LONGITUDE` | `floating point number` | Longitude of the center of the area to be covered                 |
| `TRACKER_RADIUS_DISTANCE`    | `number`                | Radius distance from the center of the area to be covered (in Km) |

## Usage

If not running the docker compose, you also need to run a postgres database yourself.

-   With docker:
    -   Build the image
    -   Use docker compose to start both containers
-   With a binary:
    -   Run `chmod +x planes-tracker`.
    -   Start the tool with `./planes-tracker`
    -   You should definetely start it in a screen or daemonize it.
-   Running from source:
    -   Start the tool with `go run .`

## Credits

-   [FlightRadar24](https://www.flightradar24.com/)
-   [ADS-B Exchange](https://adsbexchange.com/)

## Copyright

See the [license](/LICENSE).
