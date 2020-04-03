begin;

create table if not exists groups (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  defaultGroup bool,
  description text,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),
  deletionTimestamp timestamptz,

  primary key (id)
);

insert into groups (name, defaultGroup, description) values ('flatmember', true, 'Standard user account');
insert into groups (name, defaultGroup, description) values ('admin', false, 'Administrative user account, allows for access to Admin panel and API');

commit;
