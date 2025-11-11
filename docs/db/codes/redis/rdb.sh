SAVE  # 由 Redis 主进程来执行 RDB，会阻塞所有命令
# OK
BGSAVE  # 开启子进程执行 RDB，避免主进程受到影响
# Background saving started