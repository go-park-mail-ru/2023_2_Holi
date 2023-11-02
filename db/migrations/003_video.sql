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

CREATE TRIGGER modify_video_updated_at
    BEFORE UPDATE
    ON video
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

---- create above / drop below ----

DROP TABLE video;
