begin;

create table if not exists shopping_list (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  notes text,
  author text not null,
  authorLast text not null,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),
  deletionTimestamp timestamptz not null default now(),

  primary key (id),
  foreign key (author) references users(id),
  foreign key (authorLast) references users(id)
);

commit;
