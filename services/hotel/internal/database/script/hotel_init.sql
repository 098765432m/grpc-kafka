CREATE TABLE hotels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL
);

INSERT INTO hotels (id, name) VALUES
('3868a0b9-eadb-471b-8f7b-7547cc837fb2', 'Hotel Cali'),
('a312ff75-0695-4a50-bdea-4049972e99b8', 'Hotel Lisa'),
('d51d6cee-55a5-443f-be8b-82a12fe2283a', 'Hotel Fifteen');