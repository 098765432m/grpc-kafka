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
    hotel_id TEXT DEFAULT NULL
);

INSERT INTO users (id, username, password, address, email, phone_number, full_name, role, hotel_id) VALUES
('2d236bcf-15bb-43ac-a6d0-8105c14e902a','john_doe', '$2a$10$jCyE90CnRHDm4YiTN.6/beXQ5jfUUgVr.IPul0hOVyHaB38T9vktS', 'Ong Trang', 'jd@as.com', '1234567890', 'John Doe', 'GUEST', NULL),
('395901b6-5dd5-44b0-885e-859c0bfc7dee','kim_lim', '$2a$10$jCyE90CnRHDm4YiTN.6/beXQ5jfUUgVr.IPul0hOVyHaB38T9vktS', 'Ong Trang', 'kl@as.com', '1902345678', 'Kim Lim', 'GUEST', NULL);