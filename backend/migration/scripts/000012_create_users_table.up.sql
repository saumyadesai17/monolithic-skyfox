CREATE TABLE IF NOT EXISTS users 
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone  NUMERIC(10) NOT NULL UNIQUE,
    email VARCHAR(150) UNIQUE,
    name VARCHAR(100) NOT NULL,
    avatar_url TEXT,
    password_hash TEXT,
    counter_no VARCHAR(50),
    is_phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);
