begin;

create table if not exists tasks (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  notes text,
  frequency text not null,
  completed bool not null default false,
  assignee text not null,
  rotation text not null,
  rotatesBetween text not null,
  startDate int not null,
  templateId text,
  author text not null,
  authorLast text not null,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int not null default 0,

  primary key (id),
  foreign key (author) references users(id),
  foreign key (authorLast) references users(id),
  foreign key (assignee) references users(id)
);

comment on table tasks is 'The table tasks is used for tasks to store templates and active tasks.';

commit;
