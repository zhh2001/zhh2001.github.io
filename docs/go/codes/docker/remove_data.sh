# 可选步骤：确认数据不再需要，且目录没有被其他服务使用后执行
# 按正文停止相关服务并卸载软件后，再清理默认数据目录
sudo rm -rf /var/lib/docker /var/lib/containerd
