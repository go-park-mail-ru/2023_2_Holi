drop function if exists modify_video_updated_at();

CREATE TRIGGER modify_video_updated_at
    BEFORE UPDATE
    ON video
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);
