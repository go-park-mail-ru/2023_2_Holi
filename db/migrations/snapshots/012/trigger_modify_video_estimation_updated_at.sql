drop function if exists modify_video_estimation_updated_at();

CREATE TRIGGER modify_video_estimation_updated_at
    BEFORE UPDATE
    ON video_estimation
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);
