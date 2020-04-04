-- flattrack.shopping_item_to_list definition

begin;

create table if not exists shopping_item_to_list (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  itemId text references shopping_item(id) not null,
  listId text references shopping_list(id) not null,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int not null default 0,

  primary key (id)
);

comment on table shopping_item_to_list is 'The table shopping_item_to_list is used for associating shopping list items to a shopping list';

commit;
