net.setIntfPort('s1', 'h1', 1)  # s1 面向 h1 的端口号
net.setIntfPort('h1', 's1', 0)  # h1 面向 s1 的端口号
net.setIntfPort('s1', 'h2', 2)
net.setIntfPort('h2', 's1', 0)
net.setIntfPort('s1', 'h3', 3)
net.setIntfPort('h3', 's1', 0)
net.setIntfPort('s1', 'h4', 4)
net.setIntfPort('h4', 's1', 0)