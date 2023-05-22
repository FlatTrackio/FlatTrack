begin;

insert into settings
            (name, value)
values
    ('shoppingListNotes', '')
    on conflict do nothing;

commit;
