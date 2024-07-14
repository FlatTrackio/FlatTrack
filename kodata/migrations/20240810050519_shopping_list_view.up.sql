begin;

create or replace function shopping_list_view(list_id text, obtainedFilter text, orderBy text)
  returns table (shopping_list_view jsonb)
as $$
begin
  return query
    select jsonb_build_object(
      'list', (

          select row_to_json(list)
          from shopping_list list
          where id = list_id

      ),
      'templateList', (

          select row_to_json(list)
          from shopping_list list
          where id = (select templateid from shopping_list where id = list_id limit 1)

      ),
      'totalPrice', (

        with priceTimesQuantity as (
                select listid, round((price::float8 * quantity::int)::float8::numeric, 2) as priceTimesQuantity
                from shopping_item s
                where listid = list_id
        )
        select round(sum(priceTimesQuantity)::float8::numeric, 2)
        from priceTimesQuantity

      ),
      'totalPriceWithoutExcludedTags', (

          with priceTimesQuantity as (
            select listid, round((price::float8 * quantity::int)::float8::numeric, 2) as priceTimesQuantity
            from shopping_item s
            inner join shopping_list l on l.id = listid
            where l.id = list_id and tag <> any(coalesce(l.totalTagExclude, '{}')))
          select round(sum(priceTimesQuantity)::float8::numeric, 2)
          from priceTimesQuantity limit 1

      ),
      'currentPrice', (

        select total from (
          select listid, round(coalesce(sum(i.price::float8 * i.quantity::int), 0)::float8::numeric, 2) as "total"
          from shopping_item i
          inner join shopping_list s on s.id = i.listid
          where i.listid = list_id
          and not (tag = any(coalesce(s.totalTagExclude, '{}')))
          and obtained = true
          group by listid
        )

      ),
      'totalItemsObtained', (
     
        select count(*) from shopping_item s where s.listid = list_id and obtained = true

      ),
      'totalItems', (

        select count(*) from shopping_item s where s.listid = list_id

      ),
      'pricePercentage', (

        with itemsPriceObtained as (
        select listid, sum(i.price * i.quantity) as "total"
        from shopping_item i
        inner join shopping_list s on s.id = i.listid
        where i.listid = listid
              and not (tag = any(coalesce(s.totalTagExclude, '{}')))
                  and obtained = true
                  group by listid
        ),
        items as (
        select listid, sum(i.price * i.quantity) as "total"
        from shopping_item i
        inner join shopping_list s on s.id = i.listid
        where i.listid = listid
              and not (tag = any(coalesce(s.totalTagExclude, '{}')))
                  group by listid
        )
        select distinct(coalesce(round(((itemsPriceObtained.total * 100) / items.total)::float8::numeric, 0), 0)) as "percentage"
        from shopping_item
        inner join items using(listid)
        full outer join itemsPriceObtained using(listid)
        where listid = list_id

      ),
      'listTags', (

        select array_agg(row_to_json(tags))
        from (
          select distinct(tag) as "name", round(sum(price * quantity)::float8::numeric, 2) as "price" -- NOTE this is used for current price calc for each tag
          from shopping_item
          where listid = list_id
          group by tag
          order by tag
        ) tags limit 1

      ),
      'splitPrice', (

        with priceTimesQuantity as (
          select listid, (price::float8 * quantity::int) as priceTimesQuantity
          from shopping_item s
          inner join shopping_list l on l.id = listid
          where listid = list_id and l.id = listid and not (tag = any(coalesce(l.totalTagExclude, '{}'))))
        select round(sum(priceTimesQuantity)::float8::numeric / (select count(*) from users where registered = true and disabled = false), 2)
        from priceTimesQuantity limit 1

      ),
      'tags', (

        select array_agg(row_to_json(tags))
        from (
          select *
          from shopping_list_tag
          group by id, name
          order by name
        ) tags

      ),

      'items', (

        select array_agg(row_to_json(items)) from (
          select *
          from shopping_item i
          where i.listid = list_id
          and
              case
                  when (select obtainedFilter = 'true') then i.obtained = true
                  when (select obtainedFilter = 'false') then i.obtained = false
                  when (select obtainedFilter = '') then true
              end
          order by
                (case when orderBy = 'highestPrice' then price end) desc,
                (case when orderBy = 'highestPrice' then name end) asc,
                (case when orderBy = 'highestQuantity' then quantity end) desc,
                (case when orderBy = 'highestQuantity' then name end) asc,
                (case when orderBy = 'lowestPrice' then price end) asc,
                (case when orderBy = 'lowestPrice' then name end) asc,
                (case when orderBy = 'lowestQuantity' then quantity end) asc,
                (case when orderBy = 'lowestQuantity' then name end) asc,
                (case when orderBy = 'recentlyAdded' then creationTimestamp end) desc,
                (case when orderBy = 'recentlyAdded' then name end) asc,
                (case when orderBy = 'recentlyUpdated' then modificationTimestamp end) desc,
                (case when orderBy = 'recentlyUpdated' then name end) asc,
                (case when orderBy = 'lastAdded' then creationTimestamp end) asc,
                (case when orderBy = 'lastAdded' then name end) asc,
                (case when orderBy = 'lastUpdated' then modificationTimestamp end) asc,
                (case when orderBy = 'lastUpdated' then name end) asc,
                (case when orderBy = 'alphabeticalDescending' then name end) asc,
                (case when orderBy = 'alphabeticalAscending' then name end) desc,
                (case when orderBy = 'tags' then tag end) asc,
                (case when orderBy = 'tags' then name end) asc
        ) items

      )) as "shopping_list_view"

  from shopping_list s limit 1;
end;
$$ language plpgsql;

create or replace function shopping_list_view_monthly_total_price()
  returns table (shopping_list_view_monthly_total_price jsonb)
as $$
begin
  return query
    select jsonb_build_object(
      'items', (

        select array_agg(row_to_json(totalPrice)) as "totalYearMonthPrice"
        from (
                with priceTimesQuantity as (
                  select listid, round(sum(price::float8 * quantity::int)::float8::numeric, 2) as priceTimesQuantity
                  from shopping_item s
                  inner join shopping_list l on l.id = listid
                  where listid = s.listid and not (tag = any(coalesce(l.totalTagExclude, '{}')))
                  group by listid
                )
                select round(sum(priceTimesQuantity)::float8::numeric, 2), to_char(to_timestamp(creationTimestamp)::date, 'YYYY-MM') yearMonth
                from shopping_item i
                join priceTimesQuantity using(listid)
                group by yearMonth
        ) totalPrice limit 1

      )) as "shopping_list_view_total_price" ;
end;
$$ language plpgsql;

commit;
