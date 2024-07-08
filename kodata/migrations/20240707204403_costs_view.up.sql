begin;

-- TODO add item sort direction via function param
create or replace function costs_view()
  returns table (costs_view jsonb)
as $$
begin
  return query
  select jsonb_build_object(
  'totalDailyCostAverage', (
    select round((sum(amount::float8) / 4 / 7)::float8::numeric, 2)::float8 as "totalDailyCostAverage"
    from costs
    where to_timestamp(invoiceDate)::date > (to_timestamp(extract(epoch from now()))::date - '31 days'::interval)
  ),
  'totalWeeklyCostAverage', (
    select round((sum(amount::float8) / 4)::float8::numeric, 2)::float8 as "totalWeeklyCostAverage"
    from costs
    where to_timestamp(invoiceDate)::date > (to_timestamp(extract(epoch from now()))::date - '31 days'::interval)
  ),
  'totalThreeMonthAverage', (
    select round((sum(amount::float8) / 3)::float8::numeric, 2)::float8 as "totalThreeMonthAverage"
    from costs
    where to_timestamp(invoiceDate)::date > (to_timestamp(extract(epoch from now()))::date - '91 days'::interval)
  ),
  'totalYearCumulativeSpend', (
    select round((sum(amount::float8))::float8::numeric, 2)::float8 as "totalYearCumulativeSpend"
    from costs
    where date_part('year', to_timestamp(invoiceDate)::date) = date_part('year', CURRENT_DATE)
  ),
  'projectedMonthAnnualized', (
    select round((sum(amount::float8) * 12)::float8::numeric, 2)::float8 as "projectedMonthAnnualized"
    from costs
    where to_timestamp(invoiceDate)::date > (to_timestamp(extract(epoch from now()))::date - '31 days'::interval)
  ),
  'lastCostInvoiceDate', (select invoiceDate from costs order by invoiceDate desc limit 1),
  'totalPriceMonthlyCumulative', (

    select array_agg(row_to_json(monthlyPrice))
    from (
      select round(sum(amount)::float8::numeric, 2) as "total", to_char(to_timestamp(invoiceDate)::date, 'YYYY-MM') yearMonth
      from costs
      group by yearMonth
      order by yearMonth
    ) monthlyPrice

  ),
  'items', (select array_agg(row_to_json(costs) order by invoiceDate desc, title asc) from costs)) as "costs_view"
  from costs limit 1;
end;
$$ language plpgsql;

commit;
