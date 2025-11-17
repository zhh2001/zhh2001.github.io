# 卸载 Docker Engine、CLI、containerd 和 Docker Compose
sudo apt purge docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-ce-rootless-extras

# 删除所有镜像、容器和卷
sudo rm -rf /var/lib/docker
sudo rm -rf /var/lib/containerd

# 删除 source list 和 keyrings
sudo rm /etc/apt/sources.list.d/docker.sources
sudo rm /etc/apt/keyrings/docker.asc