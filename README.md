# Planes tracker

[![GitHub stars](https://img.shields.io/github/stars/LockBlock-dev/planes-tracker.svg)](https://github.com/LockBlock-dev/planes-tracker/stargazers)

Saves all the planes flying above the designated area in a sqlite3 database.

## Table of content

-   [**Installation**](#installation)
-   [**Compiling from source**](#compiling-from-source)
-   [**Configuration**](#configuration)
-   [**Usage**](#usage)
-   [**Credits**](#credits)
-   [**Copyright**](#copyright)

## Installation

-   Download [go](https://go.dev/dl/) (go 1.20 required).
-   Download or clone the project.
-   Download the binary from the [Releases](../../releases) or [build it](#compiling-from-source) yourself.
-   [Configure the tool](#configuration).

## Compiling from source

-   Use [`build.sh`](/build.sh) or use `go build` in [`cmd/planes-tracker/`](/cmd/planes-tracker/)

## Configuration

The config can be found at the root of the project.

-   Open the [`config`](/config.json) in your favorite editor.
-   Provide the latitude and longitude of the tracking zone center.
-   Then provide a radius distance (in Km) from the center.
-   You can change the poll rate (in s) of the tracker.

## Usage

-   With a binary:
    -   Run `chmod +x planes-tracker`.
    -   Start the tool with `./planes-tracker`
    -   You should definetely start it in a screen or daemonize it.
-   Running from source:
    -   Start the tool with `go run ./cmd/planes-tracker/main.go` or `cd ./cmd/planes-tracker/ && go run .`

## Credits

-   [FlightRadar24](https://www.flightradar24.com/)
-   [flightdb](https://github.com/skypies/flightdb)

## Copyright

See the [license](/LICENSE).
