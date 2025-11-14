struct eventpoll
{
    //...
    struct rb_root rbr;      // 一棵红黑树，记录要监听的 FD
    struct list_head rdlist; // 一个链表，记录就绪的 FD
    //...
};

// 1. 会在内核创建 eventpoll 结构体，返回对应的句柄 epfd
int epoll_create(int size);

// 2. 将一个 FD 添加到 epoll 的红黑树中，并设置 ep_poll_callback
// callback 触发时，就把对应的 FD 加入到 rdlist 这个就绪列表中
int epoll_ctl(
    int epfd,                 // epoll 实例的句柄
    int op,                   // 要执行的操作，包括：ADD、MOD、DEL
    int fd,                   // 要监听的 FD
    struct epoll_event *event // 要监听的事件类型：读、写、异常等
);

// 3. 检查 rdlist 列表是否为空，不为空则返回就绪的 FD 的数量
int epoll_wait(
    int epfd,                   // eventpoll 实例的句柄
    struct epoll_event *events, // 空 event 数组，用于接收就绪的 FD
    int maxevents,              // events 数组的最大长度
    int timeout                 // 超时时间，-1 不超时；0 不阻塞；大于 0 为阻塞时间
);