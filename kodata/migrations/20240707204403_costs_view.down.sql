begin;

drop materialized view if exists costs_view cascade;
drop trigger if exists tg_refresh_costs_view on costs;
drop function if exists tg_refresh_costs_view() cascade;
drop index if exists costs_id_idx;

commit;
