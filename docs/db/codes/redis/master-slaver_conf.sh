# 创建 3 个文件夹
mkdir 7001 7002 7003

# 把配置文件放进去
cp /etc/redis/redis.conf 7001
cp /etc/redis/redis.conf 7002
cp /etc/redis/redis.conf 7003

# 分别修改配置文件的端口为 7001 7002 7003
...

# 新开一个窗口，启动第 1 个 Redis
redis-server 7001/redis.conf

# 新开一个窗口，启动第 2 个 Redis
redis-server 7002/redis.conf

# 新开一个窗口，启动第 3 个 Redis
redis-server 7003/redis.conf