begin;

create table if not exists groups (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  defaultGroup bool,
  description text,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int not null default 0,

  primary key (id)
);

insert into groups (name, defaultGroup, description) values ('flatmember', true, 'Standard user account');
insert into groups (name, defaultGroup, description) values ('admin', false, 'Administrative user account, allows for access to Admin panel and API');

comment on table groups is 'The groups table is used for storing groups';

commit;
