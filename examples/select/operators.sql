SELECT
 a & b,
 x <-> y AS distance,
 (a||b)::jsonb,
 min(uppercase=>true),
 round(CAST(4 AS numeric),4)::real,
(SELECT x::json, x->>'key' AS key FROM z),
 a::character varying(12) AS trunc,
 SIN(a::double precision)::real AS redux
