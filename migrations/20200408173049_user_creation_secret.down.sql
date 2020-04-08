-- flattrack.user_creation_secret rollback definition

begin;

drop table if exists user_creation_secret;

commit;
