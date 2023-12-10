CREATE EXTENSION IF NOT EXISTS moddatetime
    WITH SCHEMA public
    CASCADE;

CREATE TABLE survey
(
    id         INT  NOT NULL,
    attribute  TEXT NOT NULL,
    rate       TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (id, attribute)
);

CREATE TRIGGER modify_user_updated_at
    BEFORE UPDATE
    ON survey
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);
