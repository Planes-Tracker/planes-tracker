const cron = require("node-cron");
const reqPlanes = require("./radar");
const config = require("./config.json");

module.exports = (radar) => {
    let previousAircrafts = [];

    cron.schedule("*/10 * * * * *", async () => {
        const you = [config.location.latitude, config.location.longitude];
        const precision = config.precision;

        const res = await reqPlanes(
            you[0] + precision,
            you[0] - precision,
            you[1] - precision,
            you[1] + precision,
            {
                FAA: true,
                SAT: true,
                MLAT: true,
                FLARM: true,
                ADSB: true,
                onGround: true,
                inAir: true,
                inactive: true,
                gliders: true,
                estimatedPositions: true,
                stats: false,
                maxAge: 14400, //4 hours
            }
        );

        if (!res.data) console.error(res.status);

        const leftAircrafts = previousAircrafts.filter((aircraft) => {
            return !res.data.aircrafts.some((a) => {
                return aircraft.id === a.id;
            });
        });

        const newAircrafts = res.data.aircrafts.filter((aircraft) => {
            return !previousAircrafts.some((a) => {
                return aircraft.id === a.id;
            });
        });

        previousAircrafts = res.data.aircrafts;

        newAircrafts.forEach((aircraft) => {
            radar.emit("plane_entered", aircraft, res.data.aircrafts);
        });

        leftAircrafts.forEach((aircraft) => {
            radar.emit("plane_left", aircraft, res.data.aircrafts);
        });
    });
};
