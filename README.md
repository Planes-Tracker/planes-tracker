# Planes tracker

Saves all the planes flying above the designated area in a sqlite3 database.

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
-   [Configure the tool](#configuration).

## Compiling from source

You can either:

-   Use the [`Makefile`](/Makefile) by running `make` in the project directory.
-   Use `go build` in the project directory.

## Configuration

The config can be found at the root of the project.

-   Open the [`config`](/config.json) in your favorite editor.
-   Provide the latitude and longitude of the tracking zone center.
-   Then provide a radius distance (in Km) from the center.
-   You can change the poll rate (in seconds) of the tracker and the database file path (default is being used in Docker).

## Config details

| Item               | Values                  | Meaning                                                           |
| ------------------ | ----------------------- | ----------------------------------------------------------------- |
| databaseName       | `text`                  | Database file path                                                |
| pollRate           | `number`                | Tracker poll rate (in seconds)                                    |
| debug              | `boolean`               | Enable debug logs                                                 |
| location.latitude  | `floating point number` | Latitude of the center of the area to be covered                  |
| location.longitude | `floating point number` | Longitude of the center of the area to be covered                 |
| radius.distance    | `number`                | Radius distance from the center of the area to be covered (in Km) |

## Usage

-   With docker:
    -   Build the image
    -   Start a container
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
