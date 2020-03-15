create table if not exists groups (
  id varchar(37) primary key not null,
  name varchar(100) not null,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),
  deletionTimestamp timestamptz not null default now()
);
