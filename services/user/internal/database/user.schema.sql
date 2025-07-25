CREATE TYPE role_enum AS ENUM (
    'GUEST',
    'ADMIN',
    'MANAGER'
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    role role_enum NOT NULL DEFAULT 'GUEST',
    hotel_id TEXT
);