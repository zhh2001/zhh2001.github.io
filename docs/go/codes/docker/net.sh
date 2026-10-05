# 1. 创建自定义 bridge 网络
sudo docker network create --driver bridge docker-demo-net

# 2. 将已有容器连接到自定义网络
sudo docker run -d --name nginx-net nginx:1.30
sudo docker network connect docker-demo-net nginx-net

# 3. 临时容器在创建时加入网络，通过容器名访问 Nginx
# 待 Nginx 启动后执行，返回其默认网页
sudo docker run --rm --network docker-demo-net alpine:3.24 \
  wget -qO- http://nginx-net:80

# 4. 查看网络和容器地址
sudo docker network ls
sudo docker network inspect docker-demo-net

# 5. 断开连接，删除练习容器和网络
sudo docker network disconnect docker-demo-net nginx-net
sudo docker stop nginx-net
sudo docker rm nginx-net
sudo docker network rm docker-demo-net
