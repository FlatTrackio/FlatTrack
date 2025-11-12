begin;

insert into system
            (name, value)
values
    ('schedulerLastRun', '{}')
    on conflict do nothing;

commit;
