CREATE TABLE video_genre
(
    video_id INTEGER REFERENCES video (id),
    genre_id INTEGER REFERENCES genre (id),
    PRIMARY KEY (video_id, genre_id)
);

---- create above / drop below ----

DROP TABLE video_genre
