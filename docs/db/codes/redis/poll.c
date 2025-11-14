// pollfd 中的事件类型
#define POLLIN   // 可读事件
#define POLLOUT  // 可写事件
#define POLLERR  // 错误事件
#define POLLNVAL // fd 未打开

// pollfd 结构
struct pollfd
{
    int fd;            // 要监听的 fd
    short int events;  // 要监听的事件类型 读、写、异常
    short int revents; // 实际发生的事件类型
};

// poll函数
int poll(
    struct pollfd *fds, // pollfd 数组，可以自定义大小
    nfds_t nfds,        // 数组元素个数
    int timeout         // 超时时间
);