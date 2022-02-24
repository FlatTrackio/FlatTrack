-- flattrack.settings definition

begin;

create table if not exists settings (
  id text default (md5(random()::text || clock_timestamp()::text)::uuid)::text not null,
  name text not null,
  value text not null,
  creationTimestamp int not null default extract('epoch' from now())::int,
  modificationTimestamp int not null default extract('epoch' from now())::int,
  deletionTimestamp int not null default 0,

  primary key (id)
);

insert into settings (name, value) values ('flatName', '');
insert into settings (name, value) values ('timezone', '');
insert into settings (name, value) values ('language', 'en_US');

comment on table settings is 'The table settings is used for admin facing settings for the FlatTrack instance';

commit;
