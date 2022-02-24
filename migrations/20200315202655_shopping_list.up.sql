begin;

create table if not exists shopping_list (
  id text default (md5(random()::text || clock_timestamp()::text)::uuid)::text not null,
  name text not null,
  notes text,
  author text not null,
  authorLast text not null,
  completed bool not null default false,
  creationTimestamp int not null default extract('epoch' from now())::int,
  modificationTimestamp int not null default extract('epoch' from now())::int,
  deletionTimestamp int not null default 0,

  primary key (id),
  foreign key (author) references users(id),
  foreign key (authorLast) references users(id)
);

comment on table shopping_list is 'The table shopping_list is used for grouping items together';

commit;
