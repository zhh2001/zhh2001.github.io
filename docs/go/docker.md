---
outline: [2, 3]
---

# Docker

本文以 **Ubuntu 24.04 LTS** 上的 **Docker Engine 29.x** 为主，使用 Linux 容器和 Bash。MySQL 示例使用 **8.4.x LTS**，Java 示例使用 **Eclipse Temurin 25 LTS**。不同版本的行为会单独注明。

Docker CLI 向 Docker daemon 发送请求，由 daemon 管理镜像、容器、存储和网络。本文使用本机的普通安装模式，命令统一加 `sudo`。如果已配置 Docker 用户组，可以省略它，但该用户组具有相当于 root 的权限。

## 1. 安装

### 1.1 安装 Docker Engine

以下步骤通过 Docker 官方 apt 仓库安装。先检查是否存在冲突的包，包括 Ubuntu 提供的 `docker.io`、`docker-buildx` 等。若已有应用依赖独立安装的 containerd 或 runc，应先安排迁移，再卸载相应包。

<<< @/go/codes/docker/uninstall_conflicting_pkg.sh

添加 apt 仓库后，选择仓库中提供的 29.x 版本，安装 Engine、CLI、containerd、Buildx 和 Compose 插件：

<<< @/go/codes/docker/install.sh

`docker-ce` 和 `docker-ce-cli` 使用同一个完整的 apt 版本字符串，其余包有各自的版本号。以后执行系统升级时仍可能更新这些包，需要固定环境时应另行管理升级策略。

首次执行 `hello-world` 会拉取测试镜像，正常输出说明 CLI 能连接 daemon 并运行容器。若服务没有启动，可执行 `sudo systemctl start docker` 后重试。

### 1.2 卸载与数据清理

卸载软件并删除本文创建的 apt 源及密钥。如果还准备清理数据，应先执行 `sudo systemctl stop docker.service docker.socket containerd.service` 停止相关服务，再卸载软件：

<<< @/go/codes/docker/uninstall_engine.sh

卸载软件不会自动清除镜像、容器、数据卷和手动修改的配置。Docker Engine 29.x 的全新安装默认使用 containerd 镜像存储，镜像内容和容器快照通常位于 `/var/lib/containerd`，数据卷等仍位于 `/var/lib/docker`。从旧版升级的环境会保留原存储后端，直到主动切换。

只有确定不再需要数据时，才执行以下清理步骤。这会删除默认目录中的容器数据和数据卷，无法通过重新安装恢复。若修改过存储目录，应核对实际路径，绑定挂载的宿主机文件也不在这两个目录的清理范围内。清理 `/var/lib/containerd` 前还要确认没有其他服务使用它。

<<< @/go/codes/docker/remove_data.sh

## 2. 快速入门

### 2.1 部署 MySQL

以下命令创建 MySQL 容器，并将数据保存在命名卷 `mysql-demo-data` 中。宿主机使用 `3307` 端口，只绑定 IPv4 回环地址，避免与本机的 `3306` 端口冲突。

<<< @/go/codes/docker/install_mysql.sh

首次启动需要初始化数据库，`docker run -d` 返回容器 ID 不表示 MySQL 已可连接。查看日志，待正式服务启动完成后，再执行 `sudo docker exec -it mysql-demo mysql -u root -p`，输入创建容器时设置的密码。

`MYSQL_ROOT_PASSWORD` 是官方镜像初始化数据库时读取的变量。已有数据库的数据卷不会因重新设置该变量而修改 root 密码。示例通过交互输入避免把密码写进命令历史，但环境变量仍保存在容器配置中，需要更严格地管理密码时可使用镜像支持的 `MYSQL_ROOT_PASSWORD_FILE`。

### 2.2 镜像和容器

**镜像（image）**包含应用运行所需的文件系统内容和运行配置，例如程序、依赖库、环境变量及默认启动命令。Linux 容器通常共享宿主机内核，镜像不需要包含一套独立的 Linux 内核。

**容器（container）**是镜像的运行实例，有自己的进程、可写层和运行配置。停止容器会保留它的可写层，删除容器则会删除这部分内容。需要独立于容器生命周期保存的数据，应放在数据卷或绑定挂载中。

**镜像注册服务（registry）**存储和分发镜像，Docker Hub 是常用的公共服务。一个 registry 中可以有多个**镜像仓库（repository）**，每个 repository 可以有多个标签和镜像版本。

`docker run` 默认先查找本地镜像，缺少时再从 registry 拉取。它不会自动搜索一个应用的所有镜像，也不会每次运行都检查远端是否更新。

### 2.3 命令解读

- `docker run`：创建并启动新容器，`-d` 表示后台运行。已有容器使用 `docker start` 启动
- `--name mysql-demo`：指定容器名，在同一个 daemon 中必须唯一，已停止但未删除的容器仍占用名字
- `-p 127.0.0.1:3307:3306`：将宿主机 IPv4 回环地址的 TCP `3307` 端口映射到容器的 TCP `3306` 端口
- `-e KEY=VALUE`：设置容器环境变量，变量是否生效取决于容器中的程序或启动脚本
- `--mount type=volume,src=mysql-demo-data,dst=/var/lib/mysql`：把命名卷挂载到容器的数据目录
- `mysql:8.4`：指定 MySQL 镜像及其标签

镜像引用常见的形式为 `[registry[:port]/][namespace/]repository[:tag]`，也可以通过 `repository@sha256:...` 指定摘要。`mysql:8.4` 是 `docker.io/library/mysql:8.4` 的简写。

未指定标签或摘要时，默认标签是 `latest`。标签是可变的名称，并不保证代表最新或稳定版本。`8.4` 限定了版本系列，但镜像维护者仍可能更新该标签指向的补丁版本或基础系统。需要严格复现时，应记录并使用镜像摘要。

没有指定宿主机 IP 的 `-p 3307:3306` 默认会发布到宿主机所有网络接口。本文的本机练习显式使用 `127.0.0.1`，若需要远程访问，再按实际需求配置监听地址和访问控制。

## 3. 常见命令

| 命令            | 用途                                                     |
| --------------- | -------------------------------------------------------- |
| `docker pull`   | 拉取镜像                                                 |
| `docker push`   | 将有写入权限的镜像推送到 registry，通常需要先登录        |
| `docker images` | 查看本地镜像，也可使用 `docker image ls`                 |
| `docker rmi`    | 删除本地镜像或标签，也可使用 `docker image rm`           |
| `docker run`    | 创建并启动新容器                                         |
| `docker stop`   | 向容器主进程发送停止信号，超时后强制结束                 |
| `docker start`  | 启动已有的已停止容器                                     |
| `docker logs`   | 查看容器收集的标准输出和标准错误，具体能力取决于日志驱动 |
| `docker exec`   | 在运行中的容器里执行一个新进程                           |
| `docker ps`     | 查看运行中的容器，`-a` 包含已停止的容器                  |
| `docker rm`     | 删除容器，通常先停止容器                                 |
| `docker build`  | 根据 Dockerfile 和构建上下文构建镜像                     |
| `docker save`   | 将镜像及其配置、标签、层导出为 tar 归档                  |
| `docker load`   | 从镜像归档加载镜像及标签                                 |

以下示例使用 `nginx:1.30`，各步按顺序执行。进入容器后的 `sh` 中执行 `exit`，返回宿主机终端后再继续。

<<< @/go/codes/docker/example.sh

`docker save` 默认不压缩归档，需要压缩时可配合 `gzip`。它保存镜像，不会备份容器可写层、数据卷或绑定挂载中的业务数据。`docker load` 支持读取相应的压缩镜像归档。

`docker logs -f nginx-demo` 会持续跟踪日志，按 `Ctrl+C` 停止跟踪，不会停止容器。`docker exec` 中使用的可执行文件必须存在于镜像里，某些精简镜像不包含 Bash，甚至不包含 shell。

普通的 `docker rm` 不能删除运行中的容器。需要强制删除时可使用 `docker rm -f`，它会用 `SIGKILL` 结束主进程，应用没有正常退出的机会。

## 4. 数据卷

**数据卷（volume）**由 Docker 管理，通过挂载出现在容器的文件系统中。它的生命周期独立于容器，适合保存数据库等持久数据。**绑定挂载（bind mount）**则直接使用 daemon 所在宿主机的指定文件或目录，适合访问配置、源码等已有文件。二者都可以绕过容器可写层，但管理方式不同。

| 命令                     | 用途                                     |
| ------------------------ | ---------------------------------------- |
| `docker volume create`   | 创建命名卷                               |
| `docker volume ls`       | 查看数据卷                               |
| `docker volume rm`       | 删除指定的、未被任何容器引用的数据卷     |
| `docker volume inspect`  | 查看数据卷配置及挂载位置等信息           |
| `docker volume prune`    | 默认清理未被任何容器引用的匿名本地卷     |
| `docker volume prune -a` | 清理未被任何容器引用的本地卷，包括命名卷 |

已停止的容器也会引用数据卷，不能仅凭“没有容器正在运行”判断卷可清理。删除卷会删除卷中的数据，使用清理命令前应核对用途。默认只清理匿名卷的行为始于 Engine 23.0 的 API 1.42，清理命名卷需要加 `-a`。

### 4.1 数据卷挂载

`-v 卷名:容器内路径` 和 `--mount type=volume,src=卷名,dst=容器内路径` 都能挂载命名卷。若卷不存在，Docker 会创建它。

把空卷挂载到容器中已有内容的目录时，Docker 默认把原目录内容复制到卷中，`volume-nocopy` 可以关闭这一步。挂载非空卷会遮住容器目标目录原有内容。

以下示例创建 Nginx 容器，修改网页后删除并重建容器，检查命名卷中的内容是否保留：

<<< @/go/codes/docker/example_v.sh

待 Nginx 启动后，可访问 `http://127.0.0.1:8081` 查看网页。练习结束后执行 `sudo docker stop nginx-volume` 和 `sudo docker rm nginx-volume`，数据卷 `nginx-demo-html` 仍会保留。只有确定不再需要这些内容时，才执行 `sudo docker volume rm nginx-demo-html`。普通 `docker rm` 不删除卷，`docker rm -v` 也只删除关联的匿名卷，不会删除这里的命名卷。

### 4.2 本地目录挂载

本地路径属于 Docker daemon 所在的宿主机。若 CLI 连接远程 daemon，不能借此直接挂载 CLI 所在机器的文件。

使用 `-v` 时，下面两种写法有不同含义：

- `-v mysql-data:/var/lib/mysql`：挂载名为 `mysql-data` 的数据卷
- `-v "$PWD/mysql/data:/var/lib/mysql"`：将当前宿主机目录下的 `mysql/data` 绑定到容器中

示例使用展开后的绝对路径，避免相对路径的处理依赖客户端版本。容器中的目标路径必须是绝对路径。`-v` 遇到不存在的宿主机路径会自动创建目录，容易掩盖路径拼写错误。`--mount type=bind` 默认要求源路径已经存在，否则报错，示例因此先创建目录。

<<< @/go/codes/docker/example_l.sh

`mysql-bind` 使用独立的数据目录和宿主机 `3308` 端口，可以与前面的 MySQL 示例同时存在。`init` 中的初始化脚本和 `conf` 中的 `.cnf` 配置以只读方式挂载。配置应放在相应配置组下，例如 `[mysqld]`，并确认容器中的 MySQL 用户能读取文件。

`/docker-entrypoint-initdb.d` 中的初始化脚本只在首次初始化数据库时执行，已有数据库不会在每次重启时重新执行。绑定一个目录会遮住镜像目标目录中原有的文件，数据目录也要保证容器进程有写入权限。

## 5. 自定义镜像

构建镜像包括准备文件系统内容和设置运行配置。Dockerfile 可以安装依赖、复制应用文件，也可以设置工作目录、默认环境变量和启动命令。

### 5.1 镜像结构

**文件系统层（layer）**记录文件的增加、修改和删除。镜像层是只读的，多个镜像可以复用相同的层，容器运行时再加上自己的可写层。Dockerfile 中的 `RUN`、`COPY`、`ADD` 等指令通常产生文件系统层，`ENV`、`CMD`、`ENTRYPOINT` 等主要修改镜像配置，不能把每一条指令都理解成一个新增的文件系统层。

**基础镜像（base image）**由 `FROM` 指定，提供当前构建阶段的起点，例如系统文件、运行时和依赖。`FROM scratch` 表示从空的起点构建，不包含现成的系统环境。

**入口（entrypoint）**和**默认命令（command）**是镜像的运行配置。`ENTRYPOINT` 常用于指定主程序，`CMD` 可以提供默认命令或默认参数，最终执行内容还受 `docker run` 参数影响。

### 5.2 Dockerfile

Dockerfile 是构建镜像的文本描述，常见指令如下：

| 指令         | 说明                                              | 示例                                     |
| ------------ | ------------------------------------------------- | ---------------------------------------- |
| `FROM`       | 指定构建阶段的基础镜像                            | `FROM eclipse-temurin:25-jre-noble`      |
| `WORKDIR`    | 设置后续指令及容器进程的工作目录                  | `WORKDIR /app`                           |
| `ENV`        | 设置保存在镜像中的环境变量                        | `ENV TZ=Asia/Shanghai`                   |
| `COPY`       | 将构建上下文中的文件复制到镜像                    | `COPY docker-demo.jar app.jar`           |
| `RUN`        | 在构建时执行命令，支持 shell 形式和 JSON 数组形式 | `RUN ["java", "-version"]`               |
| `EXPOSE`     | 声明应用预期使用的容器端口，不发布宿主机端口      | `EXPOSE 8080`                            |
| `CMD`        | 设置默认命令，或为 `ENTRYPOINT` 提供默认参数      | `CMD ["--server.port=8080"]`             |
| `ENTRYPOINT` | 设置容器的主程序                                  | `ENTRYPOINT ["java", "-jar", "app.jar"]` |

`RUN` 在构建时执行，`ENTRYPOINT` 和 `CMD` 在容器启动时使用。JSON 数组形式不会自动启动 shell，也不会自动展开 `$变量`。同时使用 exec 形式的 `ENTRYPOINT` 和 `CMD` 时，`CMD` 的内容作为默认参数追加到入口命令后，`docker run 镜像 参数...` 可以替换这些默认参数，`--entrypoint` 可以替换入口程序。

本例将 Dockerfile 命名为 `e.dockerfile`。先创建这个文本文件，再写入下面的内容。文件名中的 `e` 没有特殊含义，Docker 也不要求使用 `.dockerfile` 扩展名。默认文件名是 `Dockerfile`，使用其他名称时需要通过 `docker build -f` 指定。

下面的示例使用 Java 25 JRE，`noble` 表示镜像中的基础系统为 Ubuntu 24.04。JRE 用于运行已经构建好的应用，编译需要 JDK。准备一个包含 `Main-Class`、可在 Java 25 上运行的可执行 `docker-demo.jar`，放在 `e.dockerfile` 同一目录。这里假定应用监听容器的 `0.0.0.0:8080`，实际端口应按应用配置调整。Java 与第三方依赖的兼容性需要实际验证，不能仅根据旧 JAR 的编译版本判断。

<<< @/go/codes/docker/e.dockerfile

在包含这两个文件的目录中执行：

<<< @/go/codes/docker/build.sh

- `-f e.dockerfile`：明确指定 Dockerfile 文件，省略时默认查找构建上下文根目录中的 `Dockerfile`
- `-t docker-demo:1.0`：给构建结果设置镜像名和标签，标签由构建者定义，不指定时默认为 `latest`
- `.`：使用当前目录作为**构建上下文**，`COPY` 的源路径以此为基准，并不是指定 Dockerfile 所在目录

构建上下文应只包含需要的文件，可通过 `.dockerignore` 排除不需要发送给构建器的内容。镜像构建后，可执行 `sudo docker run -d --name java-demo -p 127.0.0.1:8082:8080 docker-demo:1.0`，待应用启动后访问 `http://127.0.0.1:8082`。`EXPOSE 8080` 不会让应用自动监听端口，也不会代替这里的 `-p`。

## 6. 网络

在本文的 Linux Docker Engine 环境中，未指定 `--network` 的普通容器默认加入 `bridge` 网络。还可以使用自定义 bridge、`host`、`none` 等模式，不能把所有容器都视为连接同一个网桥。bridge 网络适用于同一个 Docker 主机中的容器。

默认 `bridge` 网络不提供容器名的自动 DNS 解析。连接到同一个自定义 bridge 网络的容器，可以通过容器名或网络别名访问对方。同一网络内通信使用**容器端口**，不需要为了容器之间的访问设置 `-p`。容器里的 `localhost` 通常指向容器自身，需要访问别的容器时应使用对方的名字或地址。

| 命令                        | 说明                                       |
| --------------------------- | ------------------------------------------ |
| `docker network create`     | 创建网络，不指定驱动时默认为 bridge        |
| `docker network ls`         | 查看网络                                   |
| `docker network rm`         | 删除没有容器连接的指定网络                 |
| `docker network prune`      | 清理没有容器引用的自定义网络，保留内置网络 |
| `docker network connect`    | 将已有容器连接到网络                       |
| `docker network disconnect` | 将容器从网络断开                           |
| `docker network inspect`    | 查看网络配置、连接的容器及其地址           |

下面的示例先把已有 Nginx 容器加入自定义网络，再创建一个临时 Alpine 容器，用容器名发起 HTTP 请求。无需发布宿主机端口。

<<< @/go/codes/docker/net.sh

`--rm` 会在临时容器退出后自动删除该容器，命名卷不会因此删除。网络 ID、网桥名和分配的子网可能随环境变化，应以 `docker network inspect` 的结果为准。删除自定义网络前，先删除或断开连接到它的容器。
