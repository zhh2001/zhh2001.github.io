---
layout: home
title: 张恒华 | 科研与开发笔记
titleTemplate: false
head:
  - - link
    - rel: canonical
      href: https://zhh2001.github.io/
  - - meta
    - property: og:title
      content: 张恒华 | 科研与开发笔记
  - - meta
    - property: og:description
      content: 张恒华的个人网站，记录软件定义网络、P4、带内网络遥测和软件开发的学习笔记，收录个人简历与科研成果。
  - - meta
    - property: og:url
      content: https://zhh2001.github.io/
  - - meta
    - property: og:type
      content: website
  - - meta
    - name: twitter:title
      content: 张恒华 | 科研与开发笔记
  - - meta
    - name: twitter:description
      content: 张恒华的个人网站，记录软件定义网络、P4、带内网络遥测和软件开发的学习笔记，收录个人简历与科研成果。
hero:
  name: 张恒华
  text: 技术爱好者
  tagline: SDN, PDP, P4, INT, Mininet
  image:
    src: /typewriter.svg
    alt: Typewriter
  actions:
    - theme: brand
      text: Study Notes
      link: /sdn/
    - theme: alt
      text: About Me
      link: /resume
features:
  - icon:
      src: /p4/p4-logo.svg
      width: 30
    title: P4
    details: Programming Protocol-independent Packet Processors
    link: /sdn/p4
    linkText: View Notes
  - icon:
      src: /mininet/favicon.png
      width: 26
    title: Mininet
    details: An Instant Virtual Network
    link: /sdn/mininet
    linkText: View Notes
  - icon:
      src: /go/grpc.png
      width: 28
    title: gRPC
    details: A High-Performance, Open-Source Universal RPC Framework
    link: /go/grpc
    linkText: View Notes
---

---

<script setup>
import { VPTeamMembers } from 'vitepress/theme';

const members = [
  {
    avatar: '/avatar.jpg',
    name: 'Henghua Zhang',
    title: 'SDN Researcher',
    org: 'SUES',
    orgLink: 'https://www.sues.edu.cn/',
    desc: 'Focused on programmable networks, in-band network telemetry',
    links: [
      { icon: 'googlescholar', link: 'https://scholar.google.com/citations?user=nCPBFuMAAAAJ' },
      { icon: 'orcid', link: 'https://orcid.org/0009-0005-9456-8936' },
      { icon: 'github', link: 'https://github.com/zhh2001' },
      { icon: 'csdn', link: 'https://blog.csdn.net/qq_43133192' },
      { icon: 'leetcode', link: 'https://leetcode.cn/u/zhanghenghua/' },
      { icon: 'qq', link: 'mailto:1652709417@qq.com' },
    ],
  },
]
</script>

<VPTeamMembers :members />
