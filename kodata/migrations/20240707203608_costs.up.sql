begin;

create table if not exists costs (
  id text default md5(random()::text || clock_timestamp()::text)::uuid not null,
  title text not null,
  paymentType text not null,
  notes text,
  amount float8 not null,
  invoiceLink text,
  invoiceDate int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  invoicedBy text not null,
  author text not null,
  authorLast text not null,
  creationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  modificationTimestamp int not null default date_part('epoch',CURRENT_TIMESTAMP)::int,
  deletionTimestamp int not null default 0,

  primary key (id),
  foreign key (invoicedBy) references users(id),
  foreign key (author) references users(id),
  foreign key (authorLast) references users(id)
);

comment on table costs is 'The costs table is used for storing costs that flats want to keep track of';

commit;
