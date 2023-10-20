CREATE TABLE video
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100)          NOT NULL,
    description     TEXT,
    duration        VARCHAR(100)          NOT NULL,
    preview_path    VARCHAR(100)          NOT NULL,
    media_path      VARCHAR(100)          NOT NULL,
    country         VARCHAR(100),
    release_year    INTEGER
                    CONSTRAINT release_year_range
                    CHECK (release_year >= 1890
                        AND release_year <= EXTRACT(YEAR FROM CURRENT_DATE)),
    rating          FLOAT(2)
                    CONSTRAINT rating_range
                        CHECK (rating BETWEEN 0 AND 10),
    age_restriction INTEGER
                    CONSTRAINT age_restriction_range
                        CHECK (rating BETWEEN 0 AND 100),
    seasons_count   INT         DEFAULT 0 NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

---- create above / drop below ----

DROP TABLE video;
