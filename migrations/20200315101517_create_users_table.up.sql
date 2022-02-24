-- flattrack.users definition

begin;

create table if not exists users (
  id text default (md5(random()::text || clock_timestamp()::text)::uuid)::text not null,
  names text not null,
  email text not null,
  password text,
  phoneNumber text,
  birthday int,
  contractAgreement bool default false not null,
  disabled bool default false not null,
  registered bool default false not null,
  lastLogin int not null default 0,
  authNonce text default (md5(random()::text || clock_timestamp()::text)::uuid)::text not null,
  creationTimestamp int not null default extract('epoch' from now())::int,
  modificationTimestamp int not null default extract('epoch' from now())::int,
  deletionTimestamp int not null default 0,

  primary key (id)
);

comment on table users is 'The users table is used for storing user accounts';

commit;
