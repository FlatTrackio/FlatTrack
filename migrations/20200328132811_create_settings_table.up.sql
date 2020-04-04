-- flattrack.settings definition

begin;

create table if not exists settings (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  value text not null,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int,

  primary key (id)
);

insert into settings (name, value) values ('flatName', '');
insert into settings (name, value) values ('timezone', '');
insert into settings (name, value) values ('language', 'en_US');

comment on table settings is 'The table settings is used for admin facing settings for the FlatTrack instance';

commit;
