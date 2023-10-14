begin;

alter table shopping_list
add column total_tag_exclude text[];

commit;
