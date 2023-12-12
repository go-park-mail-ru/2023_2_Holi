CREATE TABLE user_subscription
(
    user_id            SERIAL PRIMARY KEY,
    subscription_start TIMESTAMPTZ DEFAULT NULL,
    subscription_up_to TIMESTAMPTZ DEFAULT NULL,
    created_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE user_subscription
    ADD CONSTRAINT fk_user_subscription_user_id
        FOREIGN KEY (user_id) REFERENCES "user" (id);

---- create above / drop below ----

DROP TABLE user_subscription