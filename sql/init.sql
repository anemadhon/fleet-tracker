-- 0001_init.sql
-- Enable PostGIS
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;

-- Vehicles table
CREATE TABLE IF NOT EXISTS vehicles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Geofences (halte)
CREATE TABLE IF NOT EXISTS geofences (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    radius DOUBLE PRECISION NOT NULL,
    geom GEOGRAPHY(POINT, 4326) GENERATED ALWAYS AS (
        ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
    ) STORED,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX geofences_geom_idx ON geofences USING GIST (geom);

-- Locations (tracking)
CREATE TABLE IF NOT EXISTS vehicle_locations (
    id SERIAL PRIMARY KEY,
    vehicle_id INT NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    speed DOUBLE PRECISION,
    heading DOUBLE PRECISION,
    created_at TIMESTAMP DEFAULT NOW(),
    geom GEOGRAPHY(POINT, 4326) GENERATED ALWAYS AS (
        ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
    ) STORED
);
CREATE INDEX locations_geom_idx ON locations USING GIST (geom);
CREATE INDEX locations_vehicle_idx ON locations (vehicle_id);

-- Geofence events (entry-only)
CREATE TABLE IF NOT EXISTS geofence_events (
    id SERIAL PRIMARY KEY,
    vehicle_id INT NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    geofence_id INT NOT NULL REFERENCES geofences(id) ON DELETE CASCADE,
    entered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    location_id INT REFERENCES locations(id),
    note TEXT
);
CREATE INDEX geofence_events_vehicle_idx ON geofence_events (vehicle_id);
CREATE INDEX geofence_events_geofence_idx ON geofence_events (geofence_id);

-- Seed (optional)
INSERT INTO geofences (name, latitude, longitude, radius)
VALUES
('Halte A', -6.2001, 106.8167, 60),
('Halte B', -6.2010, 106.8175, 60)
ON CONFLICT DO NOTHING;
