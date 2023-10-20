CREATE TABLE "user"
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100),
    email      VARCHAR(100) NOT NULL UNIQUE,
    password   VARCHAR(100) NOT NULL,
    image_path VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

---- create above / drop below ----

DROP TABLE "user";
