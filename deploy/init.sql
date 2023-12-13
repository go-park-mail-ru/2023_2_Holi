CREATE EXTENSION IF NOT EXISTS moddatetime
    WITH SCHEMA public
    CASCADE;

CREATE TABLE video
(
    id                 SERIAL PRIMARY KEY,
    name               TEXT,
    description        TEXT,
    preview_path       TEXT,
    preview_video_path TEXT                  NOT NULL,
    release_year       INTEGER
                       CONSTRAINT release_year_range
                       CHECK (release_year >= 1890
                       AND release_year <= EXTRACT(YEAR FROM CURRENT_DATE)),
    rating             FLOAT(2)
                       CONSTRAINT rating_range
                       CHECK (rating BETWEEN 0 AND 10),
    age_restriction    INTEGER
                       CONSTRAINT age_restriction_range
                       CHECK (rating BETWEEN 0 AND 100),
    seasons_count      INT         DEFAULT 0 NOT NULL,
    created_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "video" ADD COLUMN tsv tsvector;

-- UPDATE "video" SET tsv = setweight(to_tsvector(name), 'A');   это нужно добавить в скрипт после добавления данных

CREATE TRIGGER modify_video_updated_at
    BEFORE UPDATE
    ON video
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

CREATE TABLE "cast"
(
    id   serial PRIMARY KEY,
    name varchar UNIQUE NOT NULL
);

ALTER TABLE "cast" ADD COLUMN tsv tsvector;

-- UPDATE "cast" SET tsv = setweight(to_tsvector(name), 'A');   это нужно добавить в скрипт после добавления данных

CREATE INDEX ix_scenes_tsv ON "cast" USING GIN(tsv);

CREATE TABLE video_cast
(
    video_id INTEGER REFERENCES video (id),
    cast_id  INTEGER REFERENCES "cast" (id),
    PRIMARY KEY (video_id, cast_id)
);

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

CREATE TABLE video_estimation
(
    rate       INTEGER
        CONSTRAINT rate_range
            CHECK (rate BETWEEN 0 AND 10),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    video_id   INTEGER REFERENCES video (id),
--     user_id    INTEGER REFERENCES "user" (id),
    user_id    INTEGER NOT NULL,
    UNIQUE (video_id, user_id)
);

CREATE TRIGGER modify_video_estimation_updated_at
    BEFORE UPDATE
    ON video_estimation
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

CREATE TABLE genre
(
    id   serial PRIMARY KEY,
    name varchar UNIQUE NOT NULL
);

CREATE TABLE video_genre
(
    video_id INTEGER REFERENCES video (id),
    genre_id INTEGER REFERENCES genre (id),
    PRIMARY KEY (video_id, genre_id)
);

CREATE TABLE tag
(
    id   serial PRIMARY KEY,
    name varchar UNIQUE NOT NULL
);

CREATE TABLE video_tag
(
    video_id INTEGER REFERENCES video (id),
    tag_id   INTEGER REFERENCES tag (id),
    PRIMARY KEY (video_id, tag_id)
);

CREATE TABLE favourite
(
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    video_id INTEGER REFERENCES video (id),
--     user_id   INTEGER REFERENCES "user" (id),
    user_id   INTEGER NOT NULL,
    PRIMARY KEY (video_id, user_id)
);

CREATE TRIGGER modify_favourite_updated_at
    BEFORE UPDATE
    ON favourite
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);
