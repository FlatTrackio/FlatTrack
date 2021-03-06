-- flattrack.user_creation_secret definition

begin;

create table if not exists user_creation_secret (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  userId text not null,
  secret text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  valid bool not null default true,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int not null default 0,

  primary key (id),
  foreign key (userId) references users(id)
);

comment on table user_creation_secret is 'The table user_creation_secret is used for storing single use secrets for user accounts to be set up from';

commit;
