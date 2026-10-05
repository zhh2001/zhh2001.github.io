# 1. 创建命名卷并启动 Nginx
sudo docker volume create nginx-demo-html
sudo docker run -d \
  --name nginx-volume \
  -p 127.0.0.1:8081:80 \
  --mount type=volume,src=nginx-demo-html,dst=/usr/share/nginx/html \
  nginx:1.30

# 2. 查看卷并修改网页
sudo docker volume ls
sudo docker volume inspect nginx-demo-html
sudo docker exec nginx-volume sh -c 'printf "%s\n" "<h1>Docker volume demo</h1>" > /usr/share/nginx/html/index.html'

# 3. 删除容器，命名卷仍保留
sudo docker stop nginx-volume
sudo docker rm nginx-volume

# 4. 用同一个卷重建容器，查看保留的内容
sudo docker run -d \
  --name nginx-volume \
  -p 127.0.0.1:8081:80 \
  --mount type=volume,src=nginx-demo-html,dst=/usr/share/nginx/html \
  nginx:1.30
sudo docker exec nginx-volume cat /usr/share/nginx/html/index.html
