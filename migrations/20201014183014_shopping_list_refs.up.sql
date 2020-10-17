begin;

alter table shopping_list add column templateId text;
alter table shopping_item add column templateId text;

commit;
