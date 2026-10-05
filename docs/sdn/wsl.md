# WSL

WSL（Windows Subsystem for Linux）提供在 Windows 上运行 Linux 发行版和工具的环境。本仓库中的 Mininet、Scapy 和 P4 实验可以在 WSL 2 中进行，工具的安装和配置仍在相应发行版内完成。

WSL 1 通过系统调用转换层运行 Linux 程序，WSL 2 则在受管理的轻量虚拟机中运行真实 Linux 内核。涉及网络命名空间、veth 和软件交换机的实验宜使用 WSL 2，并确认当前内核支持所需功能。两种架构的区别可参考 [Microsoft 官方说明](https://github.com/MicrosoftDocs/WSL/blob/a1473a69d10c9eacdfe5e99f7496652795af28d9/WSL/compare-versions.md#comparing-wsl-1-and-wsl-2)。

本文命令以 Windows PowerShell 为主，发行版名称使用 `Ubuntu`。实际操作前，用 `wsl --list --verbose` 确认本机注册名称并替换示例中的名称。若从 WSL 的 Bash 中调用这些管理命令，使用 `wsl.exe`。命令范围可对照 [官方命令参考](https://github.com/MicrosoftDocs/WSL/blob/a1473a69d10c9eacdfe5e99f7496652795af28d9/WSL/basic-commands.md)，新选项还需查看本机 `wsl --help`。

## 安装条件

本文采用 `wsl --install` 的安装流程。[官方安装说明](https://github.com/MicrosoftDocs/WSL/blob/a1473a69d10c9eacdfe5e99f7496652795af28d9/WSL/install.md#prerequisites)要求 Windows 10 2004（Build 19041）及以上，或 Windows 11。这个条件针对下文安装流程，不能直接当作 WSL 2 架构的历史最低系统要求。更旧的系统需核对手动安装文档。

初次安装时，以管理员身份打开 PowerShell，让安装程序启用所需 Windows 功能，并按提示重启。WSL 2 还需要在 BIOS/UEFI 中启用 CPU 虚拟化，以及 Windows 的“虚拟机平台”功能。它使用 Hyper-V 架构的一部分，家庭版 Windows 也可以使用，不要求安装完整 Hyper-V 角色，参见 [官方 FAQ](https://github.com/MicrosoftDocs/WSL/blob/a1473a69d10c9eacdfe5e99f7496652795af28d9/WSL/faq.yml)。

日常进入发行版、查询状态通常不需要管理员 PowerShell。Linux 中的 `root`、`sudo` 与 Windows 管理员权限也不是同一套权限。

## 安装与更新

先查询可安装的发行版：

<<< @/sdn/codes/wsl/online.ps1

`--list` 简写为 `-l`，列表选项 `--online` 简写为 `-o`。安装时使用 `NAME` 列，`FRIENDLY NAME` 是显示名称。在线目录会变化，示例输出仅保留部分条目。需要固定 Linux 版本时，可以选择 `Ubuntu-24.04` 等带版本的名称，后续启动命令也应使用实际注册名称。

选择一个发行版安装：

<<< @/sdn/codes/wsl/install.ps1

全新环境中，裸命令 `wsl --install` 会安装 WSL 和默认 Ubuntu。已装好 WSL 后，添加发行版时显式指定 `--distribution` 或 `-d`。发行版首次启动通常会提示创建 Linux 用户和密码，这个账户独立于 Windows 账户。

WSL 2 是这个安装流程的默认架构。安装后用 `wsl --list --verbose` 确认，不能只根据发行版名称判断。

更新 WSL：

<<< @/sdn/codes/wsl/update.ps1

现代 WSL 的 `--update` 更新 WSL 软件包及随包提供的组件，包括 Linux 内核。它不负责更新 Ubuntu、Debian 内部的软件包，发行版仍需使用自己的包管理器。更新后保存工作，再按需要重新启动 WSL，使新组件生效。

## 查看状态与版本

查询已注册发行版：

<<< @/sdn/codes/wsl/list.ps1

`--verbose` 简写为列表选项 `-v`。输出中的星号表示默认发行版，`STATE` 表示运行状态，`VERSION` 的 1、2 表示该发行版采用 WSL 1 还是 WSL 2。

查询默认设置：

<<< @/sdn/codes/wsl/status.ps1

`--status` 可查看默认发行版和新发行版的默认架构，具体字段随 WSL 版本而变化。查看 WSL 软件包及组件版本则使用：

```powershell
wsl --version
```

例如，本文核验环境中的 WSL 软件包版本为 `2.6.3.0`，内核版本为 `6.6.87.2-1`，Ubuntu 的列表 `VERSION` 为 `2`。这三个值分别描述软件包、内核和运行架构，含义不同。旧的 Windows 内置 WSL 可能不识别 `--version`，应结合 `--help` 和安装方式判断。

## 启动与用户

进入默认发行版：

<<< @/sdn/codes/wsl/wsl.ps1

从 PowerShell 启动时，通常会继承 Windows 当前目录，例如 `C:\Windows` 对应 Linux 中的 `/mnt/c/Windows`。Windows Terminal 配置也可能影响起始目录。希望直接进入 Linux 用户家目录时，可以使用：

<<< @/sdn/codes/wsl/home.ps1

也可以显式指定发行版和家目录：

```powershell
wsl --distribution Ubuntu --cd ~
```

选择发行版：

<<< @/sdn/codes/wsl/distribution.ps1

`--distribution` 简写为 `-d`。指定已有 Linux 用户：

<<< @/sdn/codes/wsl/user.ps1

`--user` 简写为 `-u`，只影响本次启动，不会创建账户，也不会改变该发行版的默认用户。示例中的 root 是 Linux 用户。

设置常用发行版为默认目标：

<<< @/sdn/codes/wsl/set-default.ps1

这是选择默认发行版，与设置 WSL 1/2 架构、设置默认 Linux 用户是三种不同操作。

## 文件存放位置

Linux 工具频繁读写的代码、编译产物和实验数据，通常放在 Linux 文件系统中，例如 `~/projects`。WSL 2 从 `/mnt/c`、`/mnt/d` 访问 Windows 文件系统存在跨系统访问开销，具体影响取决于工作负载，不能笼统认为所有操作都会更慢。相关建议见 [文件系统文档](https://github.com/MicrosoftDocs/WSL/blob/a1473a69d10c9eacdfe5e99f7496652795af28d9/WSL/filesystems.md#file-storage-and-performance-across-file-systems)。

Windows 可以通过文件资源管理器的 `\\wsl.localhost\Ubuntu\` 路径访问 Linux 文件。在 Linux 当前目录执行 `explorer.exe .`，也可打开对应位置。需要同时使用 Windows 编辑器和 Linux 编译工具时，可以通过这类访问方式保持项目存放在 Linux 文件系统中。

## WSL 1/2 转换

转换已有发行版，并设置后续新发行版的默认架构：

<<< @/sdn/codes/wsl/set-version.ps1

`--set-version Ubuntu 2` 转换指定的已注册发行版，`--set-default-version 2` 只影响之后安装或导入时采用默认设置的发行版，不会批量转换已有环境。转换也不会把 Ubuntu 升级到另一个发行版本。

转换前应保存实验结果并备份重要文件。过程可能需要较长时间和额外磁盘空间，耗时与文件系统大小及存储性能有关。完成后再次用 `wsl --list --verbose` 核对架构。

## 备份与迁移

`--export` 默认把发行版根文件系统导出为 tar，`--import` 在指定位置注册一个新发行版。下面保留原来的 `Ubuntu`，将副本注册为 `Ubuntu-dev`：

<<< @/sdn/codes/wsl/export.ps1

示例要求 D 盘可用。导入名称必须与现有注册名称不同，安装位置应选择新的空目录，并留出备份文件与新文件系统所需空间。导出前先结束实验、正常停止数据库等写入服务，再终止发行版，减少应用数据处于未完成状态的可能。

导出命令会覆盖同名目标文件，因此重复备份时应更换文件名，具体行为见[导出实现](https://github.com/microsoft/WSL/blob/eb7b0cff015998c199f62b2150c1cd8b5fcfe7b3/src/windows/common/WslClient.cpp#L288-L296)。tar 备份不包含 Windows 侧的 `.wslconfig`、Windows 应用，以及 `/mnt/c` 等外部挂载目录中的文件，这些内容需要另行备份。

导入后先验证环境：

```powershell
wsl --list --verbose
wsl --distribution Ubuntu-dev --exec whoami
wsl --distribution Ubuntu-dev
```

导入的环境通常以 root 作为初始默认用户，但备份中的 `/etc/wsl.conf` 也可能已有默认用户配置。应检查实际结果。若要指定普通用户，该用户必须已经存在于导入的发行版中。在 `/etc/wsl.conf` 中添加或调整 `[user]` 段，例如：

```ini
[user]
default=alice
```

这里的 `alice` 应替换为实际 Linux 用户名。保存后，回到 PowerShell 终止 `Ubuntu-dev` 再启动，让配置生效。导入的发行版不一定有 `ubuntu.exe` 等应用启动器，因此不能直接套用启动器的 `config --default-user` 命令。详见 [导入发行版](https://github.com/MicrosoftDocs/WSL/blob/a1473a69d10c9eacdfe5e99f7496652795af28d9/WSL/use-custom-distro.md#add-wsl-specific-components-like-a-default-user)和 [用户配置](https://github.com/MicrosoftDocs/WSL/blob/a1473a69d10c9eacdfe5e99f7496652795af28d9/WSL/wsl-config.md#user-settings)。

WSL 2.6.3 还支持 `wsl --manage Ubuntu --move "D:\wsl\Ubuntu"`，可以直接移动 WSL 2 发行版。使用前应在本机帮助中确认支持，停止目标发行版并选择未被占用的位置。它不会额外生成备份，迁移前仍应保存可恢复副本。该版本的[移动实现](https://github.com/microsoft/WSL/blob/eb7b0cff015998c199f62b2150c1cd8b5fcfe7b3/src/windows/service/exe/LxssUserSession.cpp#L908-L960)要求目标发行版已经停止。

## 配置与 systemd

WSL 的两类配置文件作用范围不同，具体选项见 [官方配置说明](https://github.com/MicrosoftDocs/WSL/blob/a1473a69d10c9eacdfe5e99f7496652795af28d9/WSL/wsl-config.md)：

| 文件         | 位置与作用                                                                         |
| ------------ | ---------------------------------------------------------------------------------- |
| `.wslconfig` | Windows 用户目录 `%USERPROFILE%\.wslconfig`，配置 WSL 2 虚拟机资源及网络等全局选项 |
| `wsl.conf`   | 各发行版的 `/etc/wsl.conf`，配置该发行版的挂载、默认用户、互操作和启动行为         |

例如，限制 WSL 2 虚拟机资源时，可在 `.wslconfig` 中设置：

```ini
[wsl2]
memory=4GB
processors=2
```

这只是语法示例，应按本机容量和实验规模选择数值。资源由 WSL 2 虚拟机内的发行版共享，不是每个 Mininet Host 各有一份额度。调整 CPU、内存会改变实验条件，记录吞吐结果时也应记录这些限制。全局配置需要在 WSL 2 虚拟机停止并重新启动后生效。

使用 `systemctl` 管理 OVS 等服务前，先在 Linux 中检查 PID 1：

```shell
ps -p 1 -o comm=
```

如果已经输出 `systemd`，无需再次启用。WSL 2 的 systemd 支持要求 WSL 软件包版本至少为 `0.67.6`，发行版也需安装相应组件。需要启用时，用 `sudoedit /etc/wsl.conf` 合并以下设置到 `[boot]` 段，再停止该发行版并重新启动：

```ini
[boot]
systemd=true
```

重新启动后检查 PID 1 和服务状态：

```shell
ps -p 1 -o comm=
systemctl list-unit-files --type=service
```

具体要求见 [systemd 文档](https://github.com/MicrosoftDocs/WSL/blob/a1473a69d10c9eacdfe5e99f7496652795af28d9/WSL/systemd.md#how-to-enable-systemd)。systemd 服务存在不代表 WSL 会一直保持运行，不能仅凭开了某个服务就推断发行版的生命周期。

## 网络实验注意事项

WSL 2 默认采用 NAT 网络。Windows 通常可通过 `localhost` 访问 WSL 中监听的 TCP 服务，反方向从 Linux 访问 Windows 服务时，NAT 模式下通常使用 Windows 在虚拟网络中的地址。在 Linux 中先查看实际地址和路由：

```shell
ip -brief address
ip -4 route show default
```

默认 NAT 环境中，默认网关通常对应 Windows 在该虚拟网络中的地址。启用 DNS tunneling 后，`/etc/resolv.conf` 中的 DNS 地址可能是专用隧道地址，不能把它一律当作 Windows 服务地址。

Windows 11 22H2 及以上配合支持该功能的 WSL 版本，可在 `.wslconfig` 的 `[wsl2]` 段配置 `networkingMode=mirrored`。镜像模式支持 Linux 通过 `127.0.0.1` 访问 Windows 服务，但监听地址和 Windows、Hyper-V 防火墙规则仍会影响连通性。具体行为见 [官方网络说明](https://github.com/MicrosoftDocs/WSL/blob/a1473a69d10c9eacdfe5e99f7496652795af28d9/WSL/networking.md)。

Mininet 的 h1、h2 位于独立网络命名空间中。其拓扑地址不会因开启 WSL 的 localhost 转发或镜像网络而自动对 Windows 可达，需要另行安排路由或转发。可以先在 WSL 内按 [Mininet 笔记](./mininet.md)验证拓扑，再分别检查 Windows 与 WSL 之间的通信。

## 停止发行版与 WSL

关闭终端或执行 `exit` 结束当前会话，发行版是否随后停止还取决于剩余进程和 WSL 的空闲回收行为，不能保证它立刻停止，也不能说一定一直运行。

只停止一个发行版：

<<< @/sdn/codes/wsl/terminate.ps1

停止所有发行版和 WSL 2 虚拟机：

<<< @/sdn/codes/wsl/shutdown.ps1

这些命令会中断正在运行的实验。先保存结果、正常停止相关服务，再从 Windows PowerShell 执行。停止后可以用 `wsl --list --verbose` 查看状态，重新进入发行版会再次启动它。`--terminate` 只针对指定发行版，其他发行版仍运行时，WSL 2 虚拟机不会因此整体停止。

## 注销与卸载

`--unregister` 删除指定发行版的注册信息和根文件系统。它不是普通关机，也不会生成备份。下面以迁移得到的 `Ubuntu-dev` 为例，执行前应确认目标名称及可用备份：

<<< @/sdn/codes/wsl/unregister.ps1

对于通常的 WSL 2 安装，这会删除相应虚拟磁盘中的 Linux 数据。WSL 1 的根文件系统也会被删除，不能把这个命令理解为只删除一个 `ext4.vhdx` 文件。该操作不可用重新安装发行版来恢复原数据。

注销发行版与卸载 Windows 中的发行版应用、卸载 WSL 软件包是不同操作。需要卸载应用时，在 Windows 的应用管理中处理相应应用。
