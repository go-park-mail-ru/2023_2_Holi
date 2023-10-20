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

---- create above / drop below ----

DROP TABLE video_estimation
