CREATE TABLE episode
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    description   TEXT,
    duration      VARCHAR(100) NOT NULL,
    preview_path  VARCHAR(100) NOT NULL,
    media_path    VARCHAR(100) NOT NULL,
    number        INTEGER      NOT NULL,
    season_number INTEGER      NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    video_id      INTEGER REFERENCES video (id)
);

---- create above / drop below ----

DROP TABLE episode
