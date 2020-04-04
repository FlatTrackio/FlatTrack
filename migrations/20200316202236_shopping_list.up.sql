begin;

create table if not exists shopping_list (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  notes text,
  author text not null,
  authorLast text not null,
  completed bool not null default false,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int,

  primary key (id),
  foreign key (author) references users(id),
  foreign key (authorLast) references users(id)
);

comment on table shopping_list is 'The table shopping_list is used for shopping lists which items are associated with through shopping_item_to_list';

commit;
