-- flattrack.system definition

begin;

create table if not exists system (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  value text not null,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),

  primary key (id)
);

insert into system (name, value) values ('initialized', 'false');
insert into system (name, value) values ('jwtSecret', md5(random()::text || clock_timestamp()::text)::uuid);

commit;
