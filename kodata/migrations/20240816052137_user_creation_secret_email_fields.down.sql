begin;

alter table user_creation_secret
drop column emailSentStatus;
alter table user_creation_secret
drop column emailSentDate;

commit;
