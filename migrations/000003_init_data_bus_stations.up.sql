CREATE TABLE IF NOT EXISTS bus_stations (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);