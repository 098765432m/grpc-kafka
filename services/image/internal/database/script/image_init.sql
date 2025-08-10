CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id TEXT NOT NULL,
    format VARCHAR(10) NOT NULL,
    hotel_id UUID DEFAULT NULL
);