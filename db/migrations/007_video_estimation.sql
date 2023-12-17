CREATE TABLE video_estimation
(
    rate       INTEGER
        CONSTRAINT rate_range
            CHECK (rate BETWEEN 0 AND 10),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    video_id   INTEGER REFERENCES video (id),
    user_id    INTEGER REFERENCES "user" (id),
    UNIQUE (video_id, user_id)
);

CREATE TRIGGER modify_video_estimation_updated_at
    BEFORE UPDATE
    ON video_estimation
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

CREATE FUNCTION update_video_rating() RETURNS TRIGGER
    LANGUAGE plpgsql
AS
$$
BEGIN
    UPDATE video v
    SET rating = COALESCE(CAST((CAST(ve.rating_sum AS NUMERIC(20, 1)) / ve.rating_count) AS NUMERIC(20, 1)), 0)
    FROM (SELECT COUNT(*) AS rating_count, SUM(rate) AS rating_sum
          FROM video_estimation
          WHERE video_id = COALESCE(NEW.video_id, OLD.video_id)) ve
    WHERE v.id = COALESCE(NEW.video_id, OLD.video_id);

    RETURN NEW;
END;
$$;


CREATE OR REPLACE TRIGGER insert_update_delete_vide_estimation_rate
    AFTER INSERT OR UPDATE OF rate OR DELETE
    ON video_estimation
    FOR EACH ROW
EXECUTE FUNCTION update_video_rating();

---- create above / drop below ----

DROP TABLE video_estimation;
DROP FUNCTION update_video_rating();
