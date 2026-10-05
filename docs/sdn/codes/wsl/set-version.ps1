# 将已注册的 Ubuntu 转换为 WSL 2，转换前先备份
wsl --set-version Ubuntu 2

# 设置新发行版的默认架构，不转换已有发行版
wsl --set-default-version 2
