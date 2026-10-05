# 1. 拉取镜像并查看本地列表
sudo docker pull nginx:1.30
sudo docker images

# 2. 导出为未压缩的 tar 归档，删除本地标签后重新加载
# 此时还没有容器引用该镜像
sudo docker save -o nginx.tar nginx:1.30
sudo docker rmi nginx:1.30
sudo docker load -i nginx.tar

# 3. 创建并运行 Nginx，宿主机使用 8080 端口
sudo docker run -d --name nginx-demo -p 127.0.0.1:8080:80 nginx:1.30

# 4. 查看运行中的容器及全部容器
sudo docker ps
sudo docker ps -a
sudo docker ps --format 'table {{.ID}}\t{{.Image}}\t{{.Ports}}\t{{.Status}}\t{{.Names}}'

# 5. 停止并重新启动已有容器
sudo docker stop nginx-demo
sudo docker start nginx-demo

# 6. 查看日志
sudo docker logs nginx-demo

# 7. 进入容器，执行 exit 返回宿主机终端
sudo docker exec -it nginx-demo sh

# 8. 停止后删除容器
sudo docker stop nginx-demo
sudo docker rm nginx-demo
