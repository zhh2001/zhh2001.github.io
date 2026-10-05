# 先在 SQL 客户端查询位置和文件列表：
# SELECT @@log_bin_basename;
# SHOW BINARY LOGS;

# 以下路径替换为当前实例的实际日志文件。
mysqlbinlog /path/to/binlog.000001

# 仅用于阅读 ROW 事件，输出不是直接用于恢复的 SQL。
mysqlbinlog --base64-output=DECODE-ROWS -vv /path/to/binlog.000001
