CREATE TABLE IF NOT EXISTS SEAT_TYPES(
    id UUID PRIMARY KEY default gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    price DECIMAL(10,2) NOT NULL
)
