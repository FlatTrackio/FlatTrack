-- flattrack.settings definition

begin;

create table if not exists settings (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  value text not null,

  primary key (id)
);

insert into settings (name, value) values ('flatName', '');
insert into settings (name, value) values ('timezone', '');
insert into settings (name, value) values ('language', 'en_US');

commit;