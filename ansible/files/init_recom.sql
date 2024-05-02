CREATE EXTENSION IF NOT EXISTS moddatetime
    WITH SCHEMA public
    CASCADE;

CREATE TABLE recommendations
(
    id                      SERIAL PRIMARY KEY,
    user_id                 INT NOT NULL,
    movie_id                INT NOT NULL,
    created_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);