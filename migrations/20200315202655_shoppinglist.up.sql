create table if not exists shopping (
  id varchar(37) primary key not null,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),
  deletionTimestamp timestamptz not null default now()
);
