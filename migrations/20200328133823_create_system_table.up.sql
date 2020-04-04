-- flattrack.system definition

begin;

create table if not exists system (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  value text not null,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int not null default 0,

  primary key (id)
);

insert into system (name, value) values ('initialized', 'false');
insert into system (name, value) values ('jwtSecret', md5(random()::text || clock_timestamp()::text)::uuid);
insert into system (name, value) values ('instanceUUID', md5(random()::text || clock_timestamp()::text)::uuid);

comment on table system is 'The table system is used for managing the settings which are not managed by users or admins for the FlatTrack instance';

commit;
