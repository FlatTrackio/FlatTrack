begin;

create table if not exists shopping_items (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  price float8,
  regular bool,
  notes text,
  author text not null,
  authorLast text not null,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),
  deletionTimestamp timestamptz,

  primary key (id),
  foreign key (author) references users(id),
  foreign key (authorLast) references users(id)
);

commit;
