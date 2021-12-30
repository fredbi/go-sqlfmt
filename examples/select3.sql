WITH ls AS (select 'listed', CASE when x>0 then true else false END as flagged from lists),
v AS (select case when x is null then 'a' when x = true then 'b' else 'c' end as fv FROM values)
SELECT CASE a<b THEN a ELSE b END as minab, a, 
v.fv, CASE WHEN LOG10('listed'[0]) IS NULL THEN 'null'
WHEN LOG10('listed'[0]) = NaN THEN 'wrong' ELSE 'ok' END
    from x
