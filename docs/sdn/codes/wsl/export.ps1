# 示例使用 D 盘，先创建备份目录
New-Item -ItemType Directory -Path 'D:\backup' -Force | Out-Null

# 先在 Linux 中正常停止实验和写入服务，再停止发行版
wsl --terminate Ubuntu
if ($LASTEXITCODE -ne 0) {
    throw 'Unable to stop the source distribution.'
}

# 导出根文件系统。同名文件会被覆盖，应使用新的备份文件名
wsl --export Ubuntu 'D:\backup\ubuntu.tar'
if ($LASTEXITCODE -ne 0) {
    throw 'Export failed. Import was not started.'
}

# Ubuntu-dev 是新注册名称，第二个参数是新的安装位置
wsl --import Ubuntu-dev 'D:\wsl\ubuntu-dev' 'D:\backup\ubuntu.tar' --version 2
