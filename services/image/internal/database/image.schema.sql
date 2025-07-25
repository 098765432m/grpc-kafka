CREATE TABLE images (
    id UUID PRIMARY KEY gen-random-uuid(),
    public_id TEXT NOT NULL,
    format VARCHAR(10) NOT NULL,
    hotel_id UUID NULL
);