-- flattrack.shopping_list_tag definition

begin;

create table if not exists shopping_list_tag (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  author text not null,
  authorLast text not null,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int not null default 0,

  primary key (id)
);

comment on table shopping_list_tag is 'The table shopping_list_tag is used for storing tags to be used in shopping lists';

-- migrate old values
insert into shopping_list_tag (name, author, authorLast)
  select distinct on (tag) tag, author, author from shopping_item where tag <> '' order by tag;

commit;
