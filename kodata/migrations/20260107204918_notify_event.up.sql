begin;

create or replace function notify_event() returns trigger as $$
    declare
        data json;
        notification json;
    begin
        if (TG_OP = 'DELETE') then
            data = row_to_json(old);
        else
            data = row_to_json(new);
        end IF;
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
        PERFORM pg_notify('events',notification::text);
        return null;
    end;
$$ language plpgsql;

create or replace trigger shopping_item_notify_event
after insert or update or delete on shopping_item
    for each row execute procedure notify_event();

commit;
