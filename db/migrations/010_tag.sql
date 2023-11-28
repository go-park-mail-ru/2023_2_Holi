CREATE TABLE tag
(
    id   serial PRIMARY KEY,
    name varchar UNIQUE NOT NULL
);

---- create above / drop below ----

DROP TABLE tag
