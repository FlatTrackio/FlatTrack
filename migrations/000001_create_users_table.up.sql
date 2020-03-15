-- flattrack.users definition

begin;

create table if not exists users (
  id varchar(37) NOT NULL PRIMARY KEY,
  names varchar(100) NOT NULL,
  password text NOT NULL,
  phoneNumber varchar(100) DEFAULT NULL,
  email varchar(100) NOT NULL,
  allergies varchar(100) DEFAULT NULL,
  contractAgreement bool DEFAULT NULL,
  disabled bool DEFAULT NULL,
  groups varchar(100) DEFAULT NULL,
  hasSetPassword bool DEFAULT NULL,
  taskNotificationFrequency varchar(10) DEFAULT NULL,
  creationTimestamp timestamptz not null default now(),
  modificationTimestamp timestamptz not null default now(),
  deletionTimestamp timestamptz not null default now()
);

commit;
