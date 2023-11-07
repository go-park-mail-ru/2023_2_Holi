CREATE TABLE video_estimation
(
    rate       INTEGER
               CONSTRAINT rate_range
               CHECK (rate BETWEEN 0 AND 10),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    video_id   INTEGER REFERENCES video (id),
    user_id    INTEGER REFERENCES "user" (id),
    UNIQUE (video_id, user_id)
);

CREATE TRIGGER modify_video_estimation_updated_at
    BEFORE UPDATE
    ON video_estimation
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

---- create above / drop below ----

DROP TABLE video_estimation
