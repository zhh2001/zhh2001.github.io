# 交互输入密码，不写入 shell 历史
read -r -s -p 'MySQL root 密码: ' MYSQL_ROOT_PASSWORD
printf '\n'

sudo docker run -d \
  --name mysql-demo \
  -p 127.0.0.1:3307:3306 \
  -e TZ=Asia/Shanghai \
  -e MYSQL_ROOT_PASSWORD="$MYSQL_ROOT_PASSWORD" \
  --mount type=volume,src=mysql-demo-data,dst=/var/lib/mysql \
  mysql:8.4
unset MYSQL_ROOT_PASSWORD

# 检查初始化和正式启动日志
sudo docker logs mysql-demo
