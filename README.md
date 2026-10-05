# 张恒华的个人学习笔记

[![GitHub License](https://img.shields.io/github/license/zhh2001/zhh2001.github.io?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Website](https://img.shields.io/website?url=https://zhh2001.github.io&up_message=online&style=for-the-badge&logo=vue.js)](https://zhh2001.github.io)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/zhh2001/zhh2001.github.io/deploy.yml?style=for-the-badge&logo=github-actions)
[![GitHub deployments](https://img.shields.io/github/deployments/zhh2001/zhh2001.github.io/github-pages?style=for-the-badge&logo=githubpages)](https://github.com/zhh2001/zhh2001.github.io/actions/workflows/deploy.yml)
[![GitHub Discussions](https://img.shields.io/github/discussions/zhh2001/zhh2001.github.io?style=for-the-badge&logo=github)](https://github.com/zhh2001/zhh2001.github.io/discussions)
![GitHub Repo stars](https://img.shields.io/github/stars/zhh2001/zhh2001.github.io?style=for-the-badge&logo=refinedgithub)
![VitePress](https://img.shields.io/github/package-json/dependency-version/zhh2001/zhh2001.github.io/dev/vitepress?style=for-the-badge&logo=vitepress)

这个仓库存放我的个人学习笔记，主要记录软件定义网络与可编程数据平面的学习和实验，也整理 Go 后端开发、数据库及 LaTeX 学术写作中的常用方法。

在线阅读：<https://zhh2001.github.io/>

## 内容分类

### 科研笔记

- **网络架构与编程**：[SDN](docs/sdn/index.md)、[P4](docs/sdn/p4.md)、[P4Runtime](docs/sdn/p4runtime.md)、[INT 带内网络遥测](docs/sdn/int.md)
- **实验环境与工具**：[Mininet](docs/sdn/mininet.md)、[iPerf](docs/sdn/iperf.md)、[WSL](docs/sdn/wsl.md)
- **学术写作**：[LaTeX 排版](docs/sdn/writing/index.md)，包括公式、数值与单位、定理、图片、表格、算法、代码和参考文献

### 编程笔记

- **Go 开发**：[语言基础](docs/go/golang.md)、[Goroutine 与并发](docs/go/goroutine.md)、[Gin](docs/go/gin.md)、[gRPC](docs/go/grpc.md)、[Eino](docs/go/eino.md)
- **数据库**：[MySQL](docs/db/mysql.md)、[Redis](docs/db/redis.md)
- **容器**：[Docker](docs/go/docker.md)

## 技术栈

- [VitePress](https://vitepress.dev/) - 静态站点生成
- [MathJax](https://www.mathjax.org/) - 数学公式渲染
- [Giscus](https://giscus.app/) - 基于 GitHub Discussions 的评论系统

## 贡献

欢迎提交 Pull Request。涉及语言规则、协议行为或工具参数时，请注明适用版本并优先引用官方资料。外部引用尽量使用 DOI、RFC 或固定提交链接。代码示例应说明运行环境及其适用范围，提交前运行 `npm run docs:build` 检查文档构建。

## 联系方式

- **GitHub**：[zhh2001](https://github.com/zhh2001)
- **Email**：<1652709417@qq.com>

## 许可证

[MIT](./LICENSE) © 2023-present Henghua Zhang
