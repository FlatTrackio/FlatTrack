-- flattrack.users definition

begin;

create table if not exists users (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  names text not null,
  email text not null,
  password text,
  phoneNumber text,
  birthday int,
  contractAgreement bool default false not null,
  disabled bool default false not null,
  registered bool default false not null,
  taskNotificationFrequency int default 7 not null,
  lastLogin int not null default 0,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int not null default 0,

  primary key (id)
);

comment on table users is 'The users table is used for storing user accounts';

commit;
