#include <poll.h>

// 接口声明，struct pollfd 由头文件定义
int poll(struct pollfd *fds, nfds_t nfds, int timeout);
// timeout 单位为毫秒，负值表示无限等待，0 表示立即返回
// events 表示关注事件，revents 表示实际事件
// 常见标志包括 POLLIN、POLLOUT、POLLERR、POLLHUP、POLLNVAL
