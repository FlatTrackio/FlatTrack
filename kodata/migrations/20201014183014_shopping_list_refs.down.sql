begin;

alter table shopping_list drop column templateId;
alter table shopping_item drop column templateId;

commit;
