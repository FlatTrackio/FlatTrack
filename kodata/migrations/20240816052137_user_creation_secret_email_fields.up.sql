begin;

alter table user_creation_secret
add column emailSentStatus text not null;
alter table user_creation_secret
add column emailSentDate int not null default 0;

commit;
