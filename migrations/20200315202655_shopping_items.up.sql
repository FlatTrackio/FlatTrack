begin;

create table if not exists shopping_item (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  price float8,
  regular bool,
  notes text,
  author text not null,
  authorLast text not null,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int,

  primary key (id),
  foreign key (author) references users(id),
  foreign key (authorLast) references users(id)
);

comment on table shopping_item is 'The table shopping_item is used for shopping list items which are not necessarily associated with a list';

commit;
