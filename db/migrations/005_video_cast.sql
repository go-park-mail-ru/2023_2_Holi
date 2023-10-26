CREATE TABLE video_cast
(
    video_id INTEGER REFERENCES video (id),
    cast_id  INTEGER REFERENCES "cast" (id),
    PRIMARY KEY (video_id, cast_id)
);

---- create above / drop below ----

DROP TABLE video_cast
