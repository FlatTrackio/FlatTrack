begin;

create or replace function costs_view()
  returns table (costs_view jsonb)
as $$
begin
  return query
  select jsonb_build_object(
  'totalDailyCostAverage', (select round((sum(amount::float8) / 4 / 7)::float8::numeric, 2)::float8 as "weeklyCostAverage" from costs where to_timestamp(invoiceDate)::date > (to_timestamp(extract(epoch from now()))::date - '31 days'::interval)),
  'totalWeeklyCostAverage', (select round((sum(amount::float8) / 4)::float8::numeric, 2)::float8 as "weeklyCostAverage" from costs where to_timestamp(invoiceDate)::date > (to_timestamp(extract(epoch from now()))::date - '31 days'::interval)),
  'totalThreeMonthAverage', (select round((sum(amount::float8) / 3)::float8::numeric, 2)::float8 as "monthProjectedAnnualized" from costs where to_timestamp(invoiceDate)::date > (to_timestamp(extract(epoch from now()))::date - '91 days'::interval)),
  'totalYearCumulativeSpend', (select round((sum(amount::float8))::float8::numeric, 2)::float8 as "cumulativeYearlySpend" from costs where date_part('year', to_timestamp(invoiceDate)::date) = date_part('year', CURRENT_DATE)),
  'projectedMonthAnnualized', (select round((sum(amount::float8) * 12)::float8::numeric, 2)::float8 as "monthProjectedAnnualized" from costs where to_timestamp(invoiceDate)::date > (to_timestamp(extract(epoch from now()))::date - '31 days'::interval)),
  'lastCostInvoiceDate', (select invoiceDate from costs order by invoiceDate desc limit 1),
  'thisMonthCumulative',
    (select round((sum(amount::float8))::float8::numeric, 2)::float8
            from costs
            where to_timestamp(invoiceDate)::date >= (to_timestamp(extract(epoch from now()))::date - '4 weeks'::interval)),
  'oneMonthAgoCumulative',
    (select round((sum(amount::float8))::float8::numeric, 2)::float8
            from costs
            where to_timestamp(invoiceDate)::date >= (to_timestamp(extract(epoch from now()))::date - '8 weeks'::interval)
            and to_timestamp(invoiceDate)::date < (to_timestamp(extract(epoch from now()))::date - '4 weeks'::interval)),
  'twoMonthsAgoCumulative',
    (select round((sum(amount::float8))::float8::numeric, 2)::float8
            from costs
            where to_timestamp(invoiceDate)::date >= (to_timestamp(extract(epoch from now()))::date - '12 weeks'::interval)
            and to_timestamp(invoiceDate)::date < (to_timestamp(extract(epoch from now()))::date - '8 weeks'::interval)),
  'threeMonthsAgoCumulative',
    (select round((sum(amount::float8))::float8::numeric, 2)::float8
            from costs
            where to_timestamp(invoiceDate)::date >= (to_timestamp(extract(epoch from now()))::date - '16 weeks'::interval)
            and to_timestamp(invoiceDate)::date < (to_timestamp(extract(epoch from now()))::date - '12 weeks'::interval)),
  'items', (select array_agg(row_to_json(costs) order by invoiceDate desc) from costs)) as "costs_view"
  from costs limit 1;
end;
$$ language plpgsql;

commit;
