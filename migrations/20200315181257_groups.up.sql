begin;

create table if not exists groups (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  defaultGroup bool,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),
  deletionTimestamp timestamptz,

  primary key (id)
);

insert into groups (name, defaultGroup) values ('flatmember', true);
insert into groups (name, defaultGroup) values ('admin', false);

commit;
