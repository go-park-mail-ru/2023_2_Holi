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

CREATE TABLE user_subscription
(
    user_id            SERIAL PRIMARY KEY,
    subscription_up_to DATE DEFAULT '01-01-0001',
    created_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE user_subscription
    ADD CONSTRAINT fk_user_subscription_user_id
        FOREIGN KEY (user_id) REFERENCES "user" (id);

CREATE TRIGGER modify_user_subscription_updated_at
    BEFORE UPDATE
    ON user_subscription
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime(updated_at);
