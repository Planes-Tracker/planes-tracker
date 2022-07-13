const EventEmitter = require("events");
const emitter = require("./emitter");
const radar = new EventEmitter();
const Database = require("better-sqlite3");

const db = new Database("planes.db");

db.exec(
    "CREATE TABLE IF NOT EXISTS planes('id' VARCHAR PRIMARY KEY, 'registration' VARCHAR, 'flight' VARCHAR, 'callsign' VARCHAR, 'origin' VARCHAR, 'destination' VARCHAR, 'latitude' FLOAT, 'longitude' FLOAT, 'altitude' int, 'bearing' int, 'speed' int, 'rateOfClimb' int, 'isOnGround' BOOL, 'sqawkCode' VARCHAR, 'model' VARCHAR, 'modeSCode' VARCHAR, 'radar' VARCHAR, 'isGlider' BOOL, 'enteredAt' TIMESTAMP, 'leftAt' TIMESTAMP);"
);

const enter = db.prepare(
    "INSERT INTO planes VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
);

const left = db.prepare("UPDATE planes SET leftAt = ? WHERE id == ?;");

radar.on("plane_entered", (p, currentPlanes) => {
    const info = enter.run(
        p.id,
        p.registration,
        p.flight,
        p.callsign,
        p.origin,
        p.destination,
        p.latitude,
        p.longitude,
        p.altitude,
        p.bearing,
        p.speed,
        p.rateOfClimb,
        p.isOnGround ? 1 : 0,
        p.sqawkCode,
        p.model,
        p.modeSCode,
        p.radar,
        p.isGlider ? 1 : 0,
        p.timestamp,
        0
    );
});

radar.on("plane_left", (p, currentPlanes) => {
    const info = left.run(p.timestamp, p.id);
});

emitter(radar);

console.log("Checking for planes...");
