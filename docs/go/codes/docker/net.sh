docker network create zhh
# 066dbc62c39cb6558abbcd93e1f6f69421ff5e1e2c9c844a9c4b8a5e4d2767b7

docker network ls
# NETWORK ID     NAME      DRIVER    SCOPE
# c26ce9ed4d92   bridge    bridge    local
# 3ed15531a5cc   host      host      local
# a4486c6487a5   none      null      local
# 066dbc62c39c   zhh       bridge    local

ip a
# 1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
#     link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
#     inet 127.0.0.1/8 scope host lo
#        valid_lft forever preferred_lft forever
#     inet6 ::1/128 scope host noprefixroute 
#        valid_lft forever preferred_lft forever
# 2: enp0s3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
#     link/ether 08:00:27:70:a0:55 brd ff:ff:ff:ff:ff:ff
#     inet 10.120.100.27/16 metric 100 brd 10.120.255.255 scope global dynamic enp0s3
#        valid_lft 14103sec preferred_lft 14103sec
#     inet6 2001:250:6009:120:a00:27ff:fe70:a055/64 scope global dynamic mngtmpaddr noprefixroute 
#        valid_lft 2591999sec preferred_lft 604799sec
#     inet6 fe80::a00:27ff:fe70:a055/64 scope link 
#        valid_lft forever preferred_lft forever
# 3: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default 
#     link/ether 06:36:d5:83:fc:32 brd ff:ff:ff:ff:ff:ff
#     inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
#        valid_lft forever preferred_lft forever
#     inet6 fe80::436:d5ff:fe83:fc32/64 scope link 
#        valid_lft forever preferred_lft forever
# 12: br-066dbc62c39c: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default 
#     link/ether da:2b:f3:93:f6:08 brd ff:ff:ff:ff:ff:ff
#     inet 172.18.0.1/16 brd 172.18.255.255 scope global br-066dbc62c39c
#        valid_lft forever preferred_lft forever

docker network connect zhh mysql

# 在容器创建时就加入 zhh 网络
docker run -d --name dd -p 8080:8080 --network zhh docker-demo