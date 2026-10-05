-- x$ 视图保留数值耗时，避免对带单位的字符串进行排序。
SELECT db, query, exec_count,
       ROUND(total_latency / 1000000000000, 6) AS total_seconds,
       ROUND(avg_latency / 1000000000000, 6) AS avg_seconds,
       rows_examined, rows_sent
FROM sys.x$statement_analysis
ORDER BY total_latency DESC
LIMIT 10;
