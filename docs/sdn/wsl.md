# WSL

WSL 全称 Windows Subsystem for Linux，适用于 Linux 的 Windows 子系统。

## 安装

- 查看所有可安装的发行版 Linux：

<<< @/sdn/codes/wsl/online.ps1

`--list` 可以简写为 `-l`，`--online` 可以简写为 `-o`。

- 安装指定版本的 Linux：

<<< @/sdn/codes/wsl/install.ps1

## 查看

- 查看安装的 Linux：

<<< @/sdn/codes/wsl/list.ps1

`--verbose` 可以简写为 `-v`。

- 查看 WSL 状态：

<<< @/sdn/codes/wsl/status.ps1

## 启动

- 启动默认的 Linux：

<<< @/sdn/codes/wsl/wsl.ps1

- 启动默认 Linux 的同时进入用户目录：

<<< @/sdn/codes/wsl/home.ps1

- 启动指定版本的 Linux：

<<< @/sdn/codes/wsl/distribution.ps1

`--distribution` 可以简写为 `-d`。

- 以指定用户启动 Linux：

<<< @/sdn/codes/wsl/user.ps1

`--user` 可以简写为 `-u`。

## 关机

- 关闭全部 Linux：

<<< @/sdn/codes/wsl/shutdown.ps1

- 关闭指定 Linux：

<<< @/sdn/codes/wsl/terminate.ps1

## 卸载

- 注销并删除 Linux：

<<< @/sdn/codes/wsl/unregister.ps1
