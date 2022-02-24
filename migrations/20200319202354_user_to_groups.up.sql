-- flattrack.user_to_groups definition

begin;

create table if not exists user_to_groups (
  id text default (md5(random()::text || clock_timestamp()::text)::uuid)::text not null,
  userId text not null,
  groupId text not null,
  creationTimestamp int not null default extract('epoch' from now())::int,
  modificationTimestamp int not null default extract('epoch' from now())::int,
  deletionTimestamp int not null default 0,

  primary key (id),
  foreign key (userId) references users(id),
  foreign key (groupId) references groups(id)
);

comment on table user_to_groups is 'The table user_to_groups is used for assigning users to groups';

commit;

