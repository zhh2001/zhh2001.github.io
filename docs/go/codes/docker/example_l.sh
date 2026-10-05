# 使用独立的练习目录，首次运行时 data 应为空
mkdir -p mysql-bind-demo/{data,init,conf}
MYSQL_DEMO_DIR="$PWD/mysql-bind-demo"

read -r -s -p 'MySQL root 密码: ' MYSQL_ROOT_PASSWORD
printf '\n'

sudo docker run -d \
  --name mysql-bind \
  -p 127.0.0.1:3308:3306 \
  -e TZ=Asia/Shanghai \
  -e MYSQL_ROOT_PASSWORD="$MYSQL_ROOT_PASSWORD" \
  --mount "type=bind,src=$MYSQL_DEMO_DIR/data,dst=/var/lib/mysql" \
  --mount "type=bind,src=$MYSQL_DEMO_DIR/init,dst=/docker-entrypoint-initdb.d,readonly" \
  --mount "type=bind,src=$MYSQL_DEMO_DIR/conf,dst=/etc/mysql/conf.d,readonly" \
  mysql:8.4
unset MYSQL_ROOT_PASSWORD MYSQL_DEMO_DIR
