begin;

drop function if exists notify_event() cascade;
drop trigger if exists shopping_item_notify_event on shopping_item cascade;

commit;
