begin;

-- TODO remove resource version

alter table users drop column if exists resourceVersion;
alter table groups drop column if exists resourceVersion;
alter table settings drop column if exists resourceVersion;
alter table shopping_item drop column if exists resourceVersion;
alter table shopping_list drop column if exists resourceVersion;
alter table system drop column if exists resourceVersion;
alter table user_creation_secret drop column if exists resourceVersion;
alter table user_to_groups drop column if exists resourceVersion;

commit;
