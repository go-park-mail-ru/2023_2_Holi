CREATE TABLE video_tag
(
    video_id INTEGER REFERENCES video (id),
    tag_id   INTEGER REFERENCES tag (id),
    UNIQUE (video_id, tag_id)
);

---- create above / drop below ----

DROP TABLE video_tag
