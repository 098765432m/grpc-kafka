CREATE TABLE hotels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL
);

INSERT INTO hotels (id, name, address) VALUES
('3868a0b9-eadb-471b-8f7b-7547cc837fb2', 'Hotel Cali', 'California'),
('a312ff75-0695-4a50-bdea-4049972e99b8', 'Hotel Lisa', 'Paris'),
('d51d6cee-55a5-443f-be8b-82a12fe2283a', 'Hotel Fifteen', 'New York');