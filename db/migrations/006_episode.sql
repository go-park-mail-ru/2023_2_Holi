CREATE TABLE episode
(
    id            SERIAL PRIMARY KEY,
    name          TEXT     NOT NULL,
    description   TEXT,
    duration      INTERVAL NOT NULL,
    preview_path  TEXT     NOT NULL,
    media_path    TEXT     NOT NULL,
    number        INTEGER  NOT NULL,
    season_number INTEGER  NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    video_id      INTEGER REFERENCES video (id)
);

CREATE TRIGGER modify_episode_updated_at
    BEFORE UPDATE
    ON episode
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

---- create above / drop below ----

DROP TABLE episode
