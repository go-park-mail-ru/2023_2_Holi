CREATE TABLE favourite
(
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    video_id INTEGER REFERENCES video (id),
    user_id   INTEGER REFERENCES "user" (id),
    PRIMARY KEY (video_id, user_id)
);

CREATE TRIGGER modify_favourite_updated_at
    BEFORE UPDATE
    ON favourite
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

---- create above / drop below ----

DROP TABLE favourite
