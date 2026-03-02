CREATE TABLE IF NOT EXISTS SEATS(
    id UUID PRIMARY KEY default gen_random_uuid(),
    screen_id UUID NOT NULL,
    seat_type_id UUID NOT NULL,
    seat_label VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'AVAILABLE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(screen_id, seat_label),
    CONSTRAINT check_seat_status CHECK (status IN ('AVAILABLE', 'RESERVED', 'MAINTENANCE')),
    FOREIGN KEY (screen_id) REFERENCES SCREENS(id),
    FOREIGN KEY (seat_type_id) REFERENCES SEAT_TYPES(id)
)
