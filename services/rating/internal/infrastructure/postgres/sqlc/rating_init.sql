CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    score INT NOT NULL CHECK (score BETWEEN 1 AND 5),
    hotel_id UUID NOT NULL,
    user_id UUID NOT NULL,
    comment TEXT
);