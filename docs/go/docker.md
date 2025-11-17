---
outline: [2,3]
---

# Docker

## 1 安装

安装前要先卸载可能冲突的包：

<<< @/go/codes/docker/uninstall_conflicting_pkg.sh

卸载 Docker 时，存储在 `/var/lib/docker/` 目录下的镜像、容器等不会自动删除。如果要清除所有相关数据：

<<< @/go/codes/docker/uninstall_engine.sh

安装步骤：

<<< @/go/codes/docker/install.sh

## 2 快速入门

### 2.1 部署 MySQL

执行以下命令即可安装 MySQL：

<<< @/go/codes/docker/install_mysql.sh

### 2.2 镜像和容器

当我们利用 Docker 安装应用时，Docker 会自动搜索并下载应用**镜像（image）**。镜像不仅包含应用本身，还包含应用运行所需要的环境、配置、系统函数库。Docker 会在运行镜像时创建一个隔离环境，称为**容器（container）**。

**镜像仓库**：存储和管理镜像的平台，Docker 官方维护了一个公共仓库：Docker Hub。

### 2.3 命令解读

- `docker run`：创建并运行一个容器，`-d` 是让容器在后台运行
- `--name mysql`：给容器起名字，必须唯一
- `-p 3306:3306`：端口映射，`宿主机端口:容器内端口`
- `-e KEY=VALUE`：环境变量
- `mysql`：镜像名

镜像名称一般由两个部分组成：`[repository]:[tag]`

- 其中 `repository` 就是镜像名
- `tag` 是镜像版本

在没有指定 tag 时，默认是 latest，代表最新版本的镜像

## 3 常见命令

1. `docker pull`：从镜像仓库拉取镜像到本地
2. `docker push`：推送本地镜像到镜像仓库
3. `docker images`：查看本地镜像列表
4. `docker rmi`：删除本地镜像
5. `docker run`：创建并运行一个容器
6. `docker stop`：停止一个容器
7. `docker start`：启动一个容器
8. `docker logs`：查看容器日志
9. `docker exec`：在运行中的容器内部执行命令
10. `docker ps`：查看运行容器的状态
11. `docker rm`：删除容器
12. `docket build`：构建自定义镜像
13. `docker save`：保存本地镜像为压缩文件
14. `docker load`：加载镜像压缩文件

示例：

<<< @/go/codes/docker/example.sh

## 4 数据卷挂载

**数据卷（volume）**是一个虚拟目录，是**容器内目录**与**宿主机目录**之间映射的桥梁。

1. `docker volume create`：创建数据卷
2. `docker volume ls`：查看全部数据卷
3. `docker volume rm`：删除指定数据卷
4. `docker volume inspect`：查看数据卷详情
5. `docker volume prune`：清除数据卷

::: tip
- 在执行 `docker run` 命令时，使用 `-v 数据卷:容器内目录` 可以完成数据卷挂载
- 当创建容器时，如果挂载了数据卷且数据卷不存在，会自动创建数据卷
:::

示例：

<<< @/go/codes/docker/example_v.sh
