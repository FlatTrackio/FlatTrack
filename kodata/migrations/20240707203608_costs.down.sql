begin;

drop table if exists costs;

delete from settings where name = 'shoppingListNotes';

commit;
