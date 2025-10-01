-- Enable unaccent
CREATE EXTENSION IF NOT EXISTS unaccent;

CREATE TABLE hotels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL
);

INSERT INTO hotels (id, name, address) VALUES
('3868a0b9-eadb-471b-8f7b-7547cc837fb2', 'Hotel Cali', 'Hồ Chí Minh'),
('a312ff75-0695-4a50-bdea-4049972e99b8', 'Hotel Lisa', 'Hồ Chí Minh'),
('d51d6cee-55a5-443f-be8b-82a12fe2283a', 'Hotel Fifteen', 'Cần Thơ');

CREATE TABLE room_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    hotel_id UUID NOT NULL,
    FOREIGN KEY (hotel_id) REFERENCES hotels(id) ON DELETE CASCADE
);

INSERT INTO room_types (id, name, price, hotel_id) VALUES
('91e67b8c-1aba-44bd-a8c3-015da7350ee5', 'Phong don', '80000', '3868a0b9-eadb-471b-8f7b-7547cc837fb2'),
('b1a9e960-caef-4da8-9b12-0b467bf74244', 'Phong doi', '120000', '3868a0b9-eadb-471b-8f7b-7547cc837fb2'),
('ba5f1d28-f156-493a-bb16-7a5c212728b2', 'Phong don', '65000', 'a312ff75-0695-4a50-bdea-4049972e99b8'),
('75a56031-674c-461d-9b3c-1fce1ad8dec2', 'Phong doi', '115000', 'a312ff75-0695-4a50-bdea-4049972e99b8');

CREATE TYPE room_status AS ENUM ('AVAILABLE', 'MAINTAINED');

CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    status room_status,
    room_type_id UUID NOT NULL,
    hotel_id UUID NOT NULL,
    FOREIGN KEY (hotel_id) REFERENCES hotels(id) ON DELETE CASCADE,
    FOREIGN KEY (room_type_id) REFERENCES room_types(id) ON DELETE CASCADE
);

INSERT INTO rooms (id, name, status, room_type_id, hotel_id) VALUES
('9754c143-fdf3-4209-a7d2-66eac786cb77', '100', 'AVAILABLE', '91e67b8c-1aba-44bd-a8c3-015da7350ee5', '3868a0b9-eadb-471b-8f7b-7547cc837fb2'),
('824110e0-517c-4cfb-a3b9-489fe19cef1d', '101', 'AVAILABLE', '91e67b8c-1aba-44bd-a8c3-015da7350ee5', '3868a0b9-eadb-471b-8f7b-7547cc837fb2'),
('5b739ee4-8ac1-468b-9e89-5207f5e801d8', '102', 'AVAILABLE', 'b1a9e960-caef-4da8-9b12-0b467bf74244', '3868a0b9-eadb-471b-8f7b-7547cc837fb2'),
('070091cd-ac8c-48c2-9218-77fff85115d8', '100', 'AVAILABLE', 'ba5f1d28-f156-493a-bb16-7a5c212728b2', 'a312ff75-0695-4a50-bdea-4049972e99b8');