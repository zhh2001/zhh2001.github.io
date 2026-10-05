# 停止所有发行版及 WSL 2 虚拟机，先保存工作
wsl --shutdown

# 查询停止后的状态，不会主动启动发行版
wsl --list --verbose
#   NAME      STATE           VERSION
# * Ubuntu    Stopped         2
