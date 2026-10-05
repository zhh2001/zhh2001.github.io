# 1. 添加 Docker 官方 apt 仓库
sudo apt update
sudo apt install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && printf '%s' "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/docker.asc
EOF

sudo apt update

# 2. 列出版本，选择一个完整的 29.x apt 版本字符串
apt-cache madison docker-ce
read -r -p 'Docker Engine apt 版本（5:29.x...）: ' DOCKER_APT_VERSION

sudo apt install \
  "docker-ce=$DOCKER_APT_VERSION" \
  "docker-ce-cli=$DOCKER_APT_VERSION" \
  containerd.io docker-buildx-plugin docker-compose-plugin
unset DOCKER_APT_VERSION

# 3. 检查服务和客户端、服务端版本，再运行测试容器
sudo systemctl status docker --no-pager
sudo docker version
sudo docker run --rm hello-world
