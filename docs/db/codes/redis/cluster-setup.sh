# 在 shell 中执行，需要六个端口及其 +10000 的 cluster bus 端口未被占用
labdir=$(mktemp -d "${TMPDIR:-/tmp}/redis-cluster.XXXXXX")
for port in 7101 7102 7103 7104 7105 7106; do
  mkdir -p "$labdir/$port"
  cat > "$labdir/$port/redis.conf" <<EOF
bind 127.0.0.1
port $port
daemonize yes
pidfile $labdir/$port/redis.pid
logfile $labdir/$port/redis.log
dir $labdir/$port
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 5000
save ""
appendonly no
EOF
  redis-server "$labdir/$port/redis.conf"
  # daemonize 返回时监听端口可能尚未就绪
  for attempt in {1..50}; do
    redis-cli -p "$port" PING >/dev/null 2>&1 && break
    sleep 0.1
  done
  redis-cli -p "$port" PING >/dev/null || exit 1
done

redis-cli --cluster create \
  127.0.0.1:7101 127.0.0.1:7102 127.0.0.1:7103 \
  127.0.0.1:7104 127.0.0.1:7105 127.0.0.1:7106 \
  --cluster-replicas 1
# 检查分配计划后输入 yes
