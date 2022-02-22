SELECT
    name,
    ope,
    TIMESTAMP_MICROS(CAST(unix_time / 1000 AS Integer)) AS time,
  (unix_time - LAG(unix_time) OVER(PARTITION BY name ORDER BY unix_time)) / 1000000 AS latency_ms
FROM
    `cloudtasks_latency.sample`
    LIMIT
    1000