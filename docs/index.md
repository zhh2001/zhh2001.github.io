---
layout: home
hero:
  name: 张恒华
  text: 技术爱好者
  tagline: SDN, OpenFlow, Mininet, PDP, P4, INT
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
    sponsor: '/sponsor',
    actionText: 'Sponsor',
    links: [
      { icon: 'csdn', link: 'https://blog.csdn.net/qq_43133192' },
      { icon: 'github', link: 'https://github.com/zhh2001' },
      { icon: 'leetcode', link: 'https://leetcode.cn/u/zhanghenghua/' },
      { icon: 'orcid', link: 'https://orcid.org/0009-0005-9456-8936' },
      { icon: 'qq', link: 'mailto:1652709417@qq.com' },
      { icon: 'rss', link: 'https://github.com/zhh2001.atom' },
    ],
  },
]
</script>

<div align="center">
  <a href="/resume" class="ToResume">
    <VPTeamMembers size="small" :members />
  </a>
</div>
