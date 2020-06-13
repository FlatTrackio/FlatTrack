begin;

-- TODO remove resource version

alter table users delete column if exists resourceVersion int not null default 0;
alter table groups delete column if exists resourceVersion int not null default 0;
alter table settings delete column if exists resourceVersion int not null default 0;
alter table shopping_item delete column if exists resourceVersion int not null default 0;
alter table shopping_list delete column if exists resourceVersion int not null default 0;
alter table system delete column if exists resourceVersion int not null default 0;
alter table user_creation_secret delete column if exists resourceVersion int not null default 0;
alter table user_to_groups delete column if exists resourceVersion int not null default 0;

commit;
