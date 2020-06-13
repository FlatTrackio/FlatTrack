-- flattrack resource version add

-- TODO before beta these changes should be squashed down into the original migrations

begin;

alter table users add column resourceVersion int not null default 0;
alter table groups add column resourceVersion int not null default 0;
alter table settings add column resourceVersion int not null default 0;
alter table shopping_item add column resourceVersion int not null default 0;
alter table shopping_list add column resourceVersion int not null default 0;
alter table system add column resourceVersion int not null default 0;
alter table user_creation_secret add column resourceVersion int not null default 0;
alter table user_to_groups add column resourceVersion int not null default 0;

commit;
