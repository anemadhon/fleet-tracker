INSERT INTO bus_stations (name, latitude, longitude) 
VALUES 
    ('Kalideres', -6.2088, 106.8456),
    ('Damai', -6.2500, 106.9000),
    ('Grogol Reformasi', -6.1750, 106.8275),
    ('Kota Bambu', -6.2200, 106.8600)
ON CONFLICT DO NOTHING;