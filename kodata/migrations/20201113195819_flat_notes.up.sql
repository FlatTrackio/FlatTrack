begin;

insert into settings
            (name, value)
values
    ('flatNotes', '')
    on conflict do nothing;

commit;
