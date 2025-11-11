begin;

insert into settings
            (name, value)
values
    ('shoppingListKeepPolicy', 'Always')
    on conflict do nothing;

commit;
