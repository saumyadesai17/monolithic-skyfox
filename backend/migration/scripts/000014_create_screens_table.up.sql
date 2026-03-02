CREATE TABLE IF NOT EXISTS SCREENS(
    id UUID PRIMARY KEY default gen_random_uuid(),
    theatre_id UUID,
    name VARCHAR(50) not null,
    total_seats INT NOT NULL,
    FOREIGN KEY (theatre_id) REFERENCES THEATRES(id)
)
