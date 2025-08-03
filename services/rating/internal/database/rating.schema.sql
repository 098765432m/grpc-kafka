CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    hotel_id UUID NOT NULL,
    user_id UUID NOT NULL,
    comment TEXT
);