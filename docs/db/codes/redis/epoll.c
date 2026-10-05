#include <sys/epoll.h>

// 接口声明，struct epoll_event 由头文件定义
int epoll_create1(int flags);
int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event);
int epoll_wait(int epfd, struct epoll_event *events,
               int maxevents, int timeout);
// op 可取 EPOLL_CTL_ADD、EPOLL_CTL_MOD、EPOLL_CTL_DEL
// epoll_wait 的 timeout 单位为毫秒，-1 无限等待，0 立即返回
