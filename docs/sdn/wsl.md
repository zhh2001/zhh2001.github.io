# WSL

做 SDN 实验时绕不开 Linux，但又不想动 Windows 上现有的工具链——这种时候 WSL 基本就是性价比最高的选择。它的全名是 Windows Subsystem for Linux，可以理解成一层让 Linux 二进制直接跑在 Windows 之上的兼容层。Mininet、Scapy、P4 编译器我大多都是丢进 WSL 里跑的，本机几乎零污染。

目前用得最多的是 WSL 2，底层换成了真正的 Linux 内核（跑在轻量 Hyper-V 里），文件系统、网络栈都比 WSL 1 完整。除非有需要直接访问 Windows 路径性能的特殊场景，新装基本就直奔 2 了。

## 装之前先确认几件事

- Windows 10 21H2 / Windows 11 以上版本，老系统的 WSL 体验差很多。
- 用**管理员身份**打开 PowerShell，否则后面 `--install` 会卡在权限上。
- BIOS 里把虚拟化（Intel VT-x / AMD-V）打开，不然 WSL 2 起不来。
- "虚拟机平台"功能默认会被 `wsl --install` 自动启用，但如果之前手动关过，需要重新勾上。

## 安装

先看微软仓库里列出的可装发行版：

<<< @/sdn/codes/wsl/online.ps1

`--list` 等价于 `-l`，`--online` 等价于 `-o`，敲多了自然就用简写。

挑一个装下去：

<<< @/sdn/codes/wsl/install.ps1

第一次跑 `wsl --install` 不带任何参数的话，会顺带把 WSL 本体和默认 Ubuntu 一起搞定，装完一般需要重启。

如果以后想把 WSL 内核本身升一下（注意是内核，不是发行版里的包），用：

<<< @/sdn/codes/wsl/update.ps1

## 查看

看一眼本机装了哪些：

<<< @/sdn/codes/wsl/list.ps1

`--verbose`（即 `-v`）会多打印运行状态和 WSL 版本号，排查问题时几乎只用这个形式。

想知道当前默认发行版、默认版本是几，看 status：

<<< @/sdn/codes/wsl/status.ps1

## 启动

直接进默认发行版：

<<< @/sdn/codes/wsl/wsl.ps1

刚进去落点是 Windows 当前所在路径（通常是 `/mnt/c/...`），跨文件系统访问会慢，做实验我习惯让它一进去就跳到家目录：

<<< @/sdn/codes/wsl/home.ps1

机器上同时装了好几个发行版时，指定一个进：

<<< @/sdn/codes/wsl/distribution.ps1

`--distribution` 简写 `-d`。

需要换个身份（比如以 `root` 进去调网络，或者用某个非默认账号）：

<<< @/sdn/codes/wsl/user.ps1

`--user` 简写 `-u`。

如果嫌每次都要带 `-d`，干脆把常用的那个设成默认：

<<< @/sdn/codes/wsl/set-default.ps1

## 在 WSL 1 / 2 之间切换

老发行版偶尔会停留在 WSL 1，做 P4/Mininet 时建议切到 2：

<<< @/sdn/codes/wsl/set-version.ps1

转换过程会重新打包根文件系统，体积大的话要等几分钟，期间别强制关窗口。

## 备份与迁移

C 盘吃紧或者想把环境搬到另一台机器，用 export / import 这一对：

<<< @/sdn/codes/wsl/export.ps1

导出来就是一个普通 tar 包，拷走、放云盘、再导入都可以。导入时第二个参数是新发行版在硬盘上的落盘目录，自己挑个空目录即可。

## 关机

平时直接关终端窗口并不会真正停掉发行版，后台还挂着。需要彻底放掉资源时：

<<< @/sdn/codes/wsl/shutdown.ps1

这一条会把所有发行版一起停掉，包括 WSL 2 的虚拟机进程，比较干净。

只想停其中一个：

<<< @/sdn/codes/wsl/terminate.ps1

## 卸载

环境玩坏了想从头再来，可以直接注销，对应的 ext4 vhdx 文件也会一并删除，**没有回收站**，操作前确认有备份：

<<< @/sdn/codes/wsl/unregister.ps1
