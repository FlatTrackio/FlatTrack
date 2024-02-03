begin;

-- why is this necessary
--   create type if not exists
-- would be so much easier but it errors out
--   failed to migrate database: migration failed: syntax error at or near "not" (column 16) in line 3: begin;
do $$
begin
    if not exists (select 1 from pg_type where typname = 'task_frequency') then
      create type task_frequency as enum ('once', 'daily', 'weekly', 'fortnightly', 'monthly');
    end if;
end
$$;

do $$
begin
    if not exists (select 1 from pg_type where typname = 'task_assignee_type') then
      create type task_assignee_type as enum ('self', 'random', 'next');
    end if;
end
$$;

create table if not exists tasks (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  name text not null,
  notes text,
  assignee text,
  assigneeType task_assignee_type not null,
  targetStartTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  frequency task_frequency not null,
  templateId text,
  latestInstanceId text,
  paused bool not null default false,
  author text not null,
  authorLast text not null,
  completed bool not null default false,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int not null default 0,

  primary key (id),
  foreign key (assignee) references users(id),
  foreign key (author) references users(id),
  foreign key (authorLast) references users(id)
);

comment on table shopping_list is 'The table shopping_list is used for grouping items together';

commit;
