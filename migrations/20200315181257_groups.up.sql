begin;

create table if not exists groups (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),
  deletionTimestamp timestamptz not null default now(),

  primary key (id)
);

commit;
