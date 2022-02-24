begin;

create table if not exists shopping_item (
  id text default (md5(random()::text || clock_timestamp()::text)::uuid)::text not null,
  listId text not null,
  name text not null,
  price float8 not null default 0,
  quantity int not null default 1,
  notes text,
  obtained bool not null default false,
  tag text,
  author text not null,
  authorLast text not null,
  creationTimestamp int not null default extract('epoch' from now())::int,
  modificationTimestamp int not null default extract('epoch' from now())::int,
  deletionTimestamp int not null default 0,

  primary key (id),
  foreign key (listId) references shopping_list(id),
  foreign key (author) references users(id),
  foreign key (authorLast) references users(id)
);

comment on table shopping_item is 'The table shopping_item is used for shopping list items which are not necessarily associated with a list';

commit;
