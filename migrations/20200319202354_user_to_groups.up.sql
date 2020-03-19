-- flattrack.users definition

begin;

create table if not exists user_to_groups (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  userId text not null,
  groupId text not null,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),
  deletionTimestamp timestamptz,

  primary key (id),
  foreign key (userId) references users(id),
  foreign key (groupId) references groups(id)
);

commit;

