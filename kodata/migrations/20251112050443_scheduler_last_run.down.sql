begin;

delete from settings where name = 'schedulerLastRun';

commit;
