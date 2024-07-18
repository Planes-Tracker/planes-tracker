# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [3.2.0] - 2024-07-18

### Added

-   Added rows can't contain empty strings anymore and will use NULL value instead

### Changed

-   Unique flights index now includes origin, destination and flight number
-   Configuration with env vars only
-   Prevent duplicated flight points

## [3.1.0] - 2024-07-09

### Changed

-   Migrate from SQLite to PostgreSQL
-   Prevent ADSB-Exchange plane model to be a null string

## [3.0.1] - 2024-06-29

### Changed

-   Flight origin/destination/diverted airports and internal airline flight number are now updated on the DB in a single request

## [3.0.0] - 2024-06-27

### Changed

-   Rewrite and added ADS-B Exchange

## [2.0.0] - 2023-12-26

### Added

-   Initial release (go version)

## [1.0.0] - 2022-07-13

### Added

-   Initial release (node version)
