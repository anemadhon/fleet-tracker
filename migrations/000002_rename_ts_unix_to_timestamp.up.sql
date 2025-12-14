DROP INDEX IF EXISTS idx_vehicle_locations_vehicle_ts;

ALTER TABLE vehicle_locations 
RENAME COLUMN ts_unix TO timestamp;

CREATE INDEX IF NOT EXISTS idx_vehicle_locations_vehicle_timestamp
    ON vehicle_locations (vehicle_id, timestamp DESC);

COMMENT ON COLUMN vehicle_locations.timestamp IS 'Unix timestamp from GPS device (renamed from ts_unix)';
