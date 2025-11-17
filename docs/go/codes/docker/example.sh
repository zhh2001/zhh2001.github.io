# 1. 拉取 Nginx
docker pull nginx

# 2. 查看本地列表
docker images

# 3. 保存镜像为压缩文件
docker save -o nginx.tar nginx:latest

# 4. 删除镜像
docker rmi nginx:latest

# 5. 加载镜像
docker load -i nginx.tar 

# 6. 创建并运行 Nginx 容器
docker run -d --name nginx -p 80:80 nginx

# 7. 查看容器
docker ps
docker ps -a
docker ps --format "table {{.ID}}\t{{.Image}}\t{{.Ports}}\t{{.Status}}\t{{.Names}}"

# 8. 停止容器
docker stop nginx

# 9. 再次启动容器
docker start nginx

# 10. 查看日志
docker logs nginx
docker logs -f nginx

# 11. 进入 Nginx 容器
docker exec -it nginx bash

# 12. 删除容器
docker rm nginx
docker rm nginx -f