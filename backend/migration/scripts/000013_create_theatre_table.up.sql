CREATE TABLE IF NOT EXISTS THEATRES(
    id UUID PRIMARY KEY default gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    location TEXT NOT NULL,
    number_of_screens INT NOT NULL
)
