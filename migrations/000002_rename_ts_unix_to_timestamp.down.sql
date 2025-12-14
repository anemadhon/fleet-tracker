DROP INDEX IF EXISTS idx_vehicle_locations_vehicle_timestamp;

ALTER TABLE vehicle_locations 
RENAME COLUMN timestamp TO ts_unix;

CREATE INDEX IF NOT EXISTS idx_vehicle_locations_vehicle_ts
    ON vehicle_locations (vehicle_id, ts_unix DESC);

COMMENT ON COLUMN vehicle_locations.ts_unix IS 'Unix timestamp from GPS device';