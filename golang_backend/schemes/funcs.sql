create function expire_session_delete_old_rows() returns trigger
    language plpgsql
as
$$
BEGIN
    WITH ses_id AS (
        DELETE FROM session WHERE created_at < NOW() - INTERVAL '1 minute' returning id
    )
    DELETE FROM "user" WHERE id = ses_id;
    RETURN NEW;
END;
$$;

CREATE TRIGGER expire_session_delete_old_rows_trigger
    AFTER INSERT ON session
EXECUTE PROCEDURE expire_session_delete_old_rows();

CREATE OR REPLACE FUNCTION allProd() returns setof product
    language plpgsql
as $$
begin
    return query
        select * from product;
end;
$$;