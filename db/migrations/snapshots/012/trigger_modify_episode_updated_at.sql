drop function if exists modify_episode_updated_at();

CREATE TRIGGER modify_episode_updated_at
    BEFORE UPDATE
    ON episode
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);
