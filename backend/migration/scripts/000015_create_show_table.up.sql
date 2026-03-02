CREATE EXTENSION IF NOT EXISTS "pgcrypto"; 
CREATE TABLE IF NOT EXISTS SHOW(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    movie_imdb_id VARCHAR(20) NOT NULL,
    screen_id UUID NOT NULL,
    theatre_id UUID NOT NULL,

    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,

    status VARCHAR(20) DEFAULT 'ACTIVE',

    seat_layout_snapshot JSONB,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_screen
        FOREIGN KEY(screen_id)
        REFERENCES screens(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_theatre
        FOREIGN KEY(theatre_id)
        REFERENCES theatres(id)
        ON DELETE RESTRICT
);
