# 需求：
#   - 创建 Nginx 容器，修改容器内的 html 目录下的 index.html 文件内容
#   - 将静态资源部署到 nginx 的 html 目录

# 1. 创建 nginx 容器
docker run -d --name nginx -p 80:80 -v html:/usr/share/nginx/html nginx

# 2. 查看全部数据卷
docker volume ls

# 3. 查看数据卷详情
docker volume inspect html