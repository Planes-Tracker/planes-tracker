const axios = require("axios");

module.exports = async (north, south, west, east, opt = {}) => {
    const baseURL = "https://data-live.flightradar24.com/zones/fcgi/feed.js";

    const params = new URLSearchParams({
        bounds: [north, south, west, east].join(","),
        faa: opt.FAA ? 1 : 0,
        satellite: opt.SAT ? 1 : 0,
        mlat: opt.MLAT ? 1 : 0,
        flarm: opt.FLARM ? 1 : 0,
        adsb: opt.ADSB ? 1 : 0,
        gnd: opt.onGround ? 1 : 0,
        air: opt.inAir ? 1 : 0,
        vehicles: opt.inactive ? 1 : 0,
        gliders: opt.gliders ? 1 : 0,
        estimated: opt.estimatedPositions ? 1 : 0,
        stats: opt.stats ? 1 : 0,
        maxage: opt.maxAge || 14400, //4 hours
    });

    let res = await axios({
        url: `${baseURL}?${decodeURIComponent(params)}`,
        method: "GET",
        headers: {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0",
            Accept: "application/json",
        },
    });

    if (res.status !== 200 || !res.data) return { ok: false, status: "api error" };

    const data = {
        version: res.data.version,
        fullCount: res.data.full_count,
        aircrafts: [],
    };

    if (opt.stats) data.stats = res.data.stats;

    res = Object.entries(res.data);

    for (let i = 0; i < res.length; i++) {
        const d = res[i][1];
        if (Array.isArray(d)) {
            data.aircrafts.push({
                id: res[i][0],
                registration: d[9] || null,
                flight: d[13] || null,
                callsign: d[16] || null, // ICAO ATC call signature
                origin: d[11] || null, // airport IATA code
                destination: d[12] || null, // airport IATA code
                latitude: d[1],
                longitude: d[2],
                altitude: d[4], // in feet
                bearing: d[3], // in degrees
                speed: d[5] || null, // in knots
                rateOfClimb: d[15], // ft/min
                isOnGround: !!d[14],
                squawkCode: d[6], // https://en.wikipedia.org/wiki/Transponder_(aeronautics)
                model: d[8] || null, // ICAO aircraft type designator
                modeSCode: d[0] || null, // ICAO aircraft registration number
                radar: d[7], // F24 "radar" data source ID
                isGlider: !!d[17],
                timestamp: d[10] || null,
            });
        }
    }

    return { ok: true, data: data };
};
