## 把整个发行版打包成 tar，方便备份或换机器
wsl --export Ubuntu D:\backup\ubuntu.tar

## 导入到新位置（第二个参数是落盘目录，自己起名）
wsl --import Ubuntu-dev D:\wsl\ubuntu-dev D:\backup\ubuntu.tar --version 2
