# 继续使用 cluster-setup.sh 创建的 labdir
mkdir -p "$labdir/7107"
cat > "$labdir/7107/redis.conf" <<EOF
bind 127.0.0.1
port 7107
daemonize yes
pidfile $labdir/7107/redis.pid
logfile $labdir/7107/redis.log
dir $labdir/7107
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 5000
save ""
appendonly no
EOF
redis-server "$labdir/7107/redis.conf"
for attempt in {1..50}; do
  redis-cli -p 7107 PING >/dev/null 2>&1 && break
  sleep 0.1
done
redis-cli -p 7107 PING >/dev/null || exit 1
redis-cli --cluster add-node 127.0.0.1:7107 127.0.0.1:7101
