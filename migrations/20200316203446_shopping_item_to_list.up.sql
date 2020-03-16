begin;

create table if not exists shopping_item_to_list (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  item text references shopping_items(id) not null,
  list text references shopping_list(id) not null,

  primary key (id)
);

commit;
