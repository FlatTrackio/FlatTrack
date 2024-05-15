begin;

create table if not exists board_items (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  title text not null,
  body text,
  author text not null,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int not null default 0,

  primary key (id),
  foreign key (author) references users(id)
);

comment on table board_items is 'The table board_items is used for storing entries for board';

commit;
