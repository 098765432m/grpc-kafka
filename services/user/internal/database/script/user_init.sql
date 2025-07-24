CREATE TYPE role_enum AS ENUM (
    'GUEST',
    'ADMIN',
    'MANAGER'
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    role role_enum NOT NULL DEFAULT 'GUEST',
    hotel_id TEXT
);

INSERT INTO users (username, password, email, phone_number, full_name, role, hotel_id) VALUES
('john_doe', '113446', 'jd@as.com', '1234567890', 'John Doe', 'GUEST', NULL),