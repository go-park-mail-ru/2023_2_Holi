CREATE TABLE "cast"
(
    id   serial PRIMARY KEY,
    name varchar UNIQUE NOT NULL,
    birthday TEXT,
    place TEXT,
    carier TEXT,
    imgPath TEXT
);

---- create above / drop below ----

DROP TABLE "cast"
