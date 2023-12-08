CREATE EXTENSION IF NOT EXISTS moddatetime
    WITH SCHEMA public
    CASCADE;

CREATE TABLE "user"
(
    id         SERIAL PRIMARY KEY,
    name       TEXT,
    email      TEXT NOT NULL UNIQUE,
    password   BYTEA NOT NULL UNIQUE,
    image_path TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER modify_user_updated_at
    BEFORE UPDATE
    ON "user"
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);
