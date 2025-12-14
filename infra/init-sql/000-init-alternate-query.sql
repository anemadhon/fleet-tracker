CREATE TABLE IF NOT EXISTS vehicle_locations (
    id          BIGSERIAL PRIMARY KEY,
    vehicle_id  VARCHAR(50) NOT NULL,
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    timestamp   BIGINT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vehicle_locations_vehicle_ts
    ON vehicle_locations (vehicle_id, timestamp DESC);

CREATE TABLE IF NOT EXISTS bus_stations (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO bus_stations (name, latitude, longitude) 
VALUES 
    ('Kalideres', -6.2088, 106.8456),
    ('Damai', -6.2500, 106.9000),
    ('Grogol Reformasi', -6.1750, 106.8275),
    ('Kota Bambu', -6.2200, 106.8600)
ON CONFLICT DO NOTHING;