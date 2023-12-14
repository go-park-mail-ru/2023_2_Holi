CREATE OR REPLACE PROCEDURE update_video_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE video v
    SET rating = (ve.ratining_sum / ve.rating_count)
    FROM (SELECT video_id, count(*) as rating_count, sum(rating) as rating_sum
        FROM video_estimation) ve
    WHERE v.id = ve.video_id
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER insert_or_update_vide_estimation_rate
    AFTER INSERT OR UPDATE
    ON video_estimation
    FOR rate
EXECUTE PROCEDURE after_insert_or_update();


---- create above / drop below ----


