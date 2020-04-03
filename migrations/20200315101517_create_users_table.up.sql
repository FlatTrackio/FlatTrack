-- flattrack.users definition

begin;

create table if not exists users (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  names text not null,
  email text not null,
  password text,
  phoneNumber text,
  birthday timestamptz,
  contractAgreement bool default false not null,
  disabled bool default false not null,
  registered bool default false not null,
  taskNotificationFrequency int default 7 not null,
  lastLogin timestamptz,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),
  deletionTimestamp timestamptz,

  primary key (id)
);

commit;
