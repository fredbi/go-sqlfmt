SELECT
  ABS(a), AVG(b+sum(x+GREATEST(y,z))) AS z,
  LEAST(STDDEV(SIN(y*SQRT(z))),STDDEV_SAMP(y)),
  ST_MakeEnvelope(minX, minY,maxX,maxY)
FROM
  numbers
