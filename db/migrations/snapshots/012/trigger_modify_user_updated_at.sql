drop function if exists modify_user_updated_at();

CREATE TRIGGER modify_user_updated_at
    BEFORE UPDATE
    ON "user"
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);
