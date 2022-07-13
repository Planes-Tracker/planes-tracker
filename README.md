# Planes tracker

[![axios](https://img.shields.io/github/package-json/dependency-version/LockBlock-dev/planes-tracker/axios)](https://www.npmjs.com/package/axios) [![node-cron](https://img.shields.io/github/package-json/dependency-version/LockBlock-dev/planes-tracker/node-cron)](https://www.npmjs.com/package/node-cron) [![better-sqlite3](https://img.shields.io/github/package-json/dependency-version/LockBlock-dev/planes-tracker/better-sqlite3)](https://www.npmjs.com/package/better-sqlite3)

[![GitHub stars](https://img.shields.io/github/stars/LockBlock-dev/planes-tracker.svg)](https://github.com/LockBlock-dev/planes-tracker/stargazers)

Saves all the planes flying above the designated area in a sqlite3 database.

## Installation

-   Install [NodeJS](https://nodejs.org).
-   Download or clone the project.
-   Run `npm install`.
-   In the [config.json](./config.json), you need to edit the location with your latitude and longitude:

```json
{
    "location": {
        "latitude": 0,
        "longitude": 0
    },
    "precision": 0.1
}
```

-   Run `node index.js` OR `npm start`.

## Credits

-   [FlightRadar24](https://www.flightradar24.com/)
-   [flightradar24-client](https://www.npmjs.com/package/flightradar24-client)

## Copyright

See the [license](/LICENSE)
