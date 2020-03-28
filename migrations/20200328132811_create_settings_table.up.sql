-- flattrack.settings definition

begin;

create table if not exists settings (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  value text not null,

  primary key (id)
);

commit;
