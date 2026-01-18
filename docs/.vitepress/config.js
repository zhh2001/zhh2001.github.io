import { defineConfig } from 'vitepress'

const SDNNoteItems = [
  { text: 'P4 Language', link: '/sdn/p4' },
  { text: 'P4 Exercise', link: '/sdn/p4_exercise' },
  { text: 'P4 INT', link: '/sdn/int' },
  { text: 'P4 Utils', link: '/sdn/p4_utils' },
  { text: 'WSL', link: '/sdn/wsl' },
  { text: 'Mininet', link: '/sdn/mininet' },
  { text: 'iPerf', link: '/sdn/iperf' },
  {
    text: 'Writing', link: '/sdn/writing/', items: [
      { text: 'Formula', link: '/sdn/writing/formula' },
      { text: 'Algorithm', link: '/sdn/writing/algorithm' },
      { text: 'Table', link: '/sdn/writing/table' },
      { text: 'Bibliography', link: '/sdn/writing/bibliography' }
    ]
  }
]

const ResearchNotes = [
  { text: 'SDN', link: '/sdn/', items: SDNNoteItems },
]
const TechNotes = [
  {
    text: 'Go', items: [
      { text: 'GoLang', link: '/go/golang' },
      { text: 'Goroutine', link: '/go/goroutine' },
      { text: 'Gin', link: '/go/gin' },
      { text: 'gRPC', link: '/go/grpc' },
      { text: 'Eino', link: '/go/eino' },
      { text: 'Docker', link: '/go/docker' }
    ]
  },
  {
    text: 'Database',
    items: [
      { text: 'MySQL', link: '/db/mysql' },
      { text: 'Redis', link: '/db/redis' },
    ]
  },
  {
    text: 'Interview',
    items: [
      { text: 'GoLang', link: '/interview/go' },
      { text: 'MySQL', link: '/interview/mysql' },
      { text: 'Redis', link: '/interview/redis' },
    ]
  }
]

const keywords = [
  '张恒华',
  '软件定义网络', '可编程数据平面', '带内网络遥测',
  'Henghua Zhang',
  'SDN', 'Software Defined Network', 'Software-Defined Networking',
  'Programmable Data Plane', 'PDP',
  'Programming Protocol-independent Packet Processors', 'P4', 'P4 Lang', 'P4 Language',
  'P4Runtime', 'P4Runtime API',
  'V1Model',
  'BMv2', 'Behavioral Model version 2',
  'INT', 'In-band Network Telemetry',
  'OpenFlow',
  'Mininet',
  'iPerf',
  'Scapy',
  'Go', 'GoLang',
  'RPC', 'gRPC',
  'Goroutine', 'GMP',
  'Gin',
  'Docker',
  'Eino',
  'Database',
  'Redis',
  'MySQL',
  'Protocol Buffer', 'protobuf'
]

const description_cn = '张恒华的个人网站，存放学习笔记与个人简历。研究领域为软件定义网络、可编程数据平面。'
const description_en = 'This is Henghua Zhang\'s personal website, where study notes and a resume are stored. The research focus is on SDN (Software-Defined Networking) and PDP (Programmable Data Plane).'

export default defineConfig({
  title: '张恒华',
  lang: 'zh-Hans',
  description: description_cn + description_en,
  appearance: true,
  cleanUrls: true,
  head: [
    ['meta', { name: 'keywords', content: keywords.join(', ') }],
    ['meta', { name: 'google-site-verification', content: 'wMOTcBwCiCMV7ESftQRY3Glvq8UL4xzUKrZ-1wjOpqM' }],
    ['link', { rel: 'icon', href: '/favicon.svg' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.googleapis.com' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' }],
    ['link', { href: 'https://fonts.googleapis.com/css2?family=Noto+Serif+SC:wght@200..900&display=swap', rel: 'stylesheet' }],
    ['script', { async: '', src: 'https://www.googletagmanager.com/gtag/js?id=G-414XCQ3MDV' }],
    ['script', {},
      `
      window.dataLayer = window.dataLayer || [];
      function gtag() {
        dataLayer.push(arguments);
      }
      gtag('js', new Date());
      gtag('config', 'G-414XCQ3MDV');
      `
    ]
  ],
  markdown: {
    math: true,
    languageAlias: {
      'p4': 'c++',
      'pseudo': 'c++',
      'conf': 'sh',
      'cnf': 'ini',
    }
  },
  themeConfig: {
    siteTitle: '张恒华',
    editLink: {
      pattern: ({ filePath }) => `https://github.com/zhh2001/zhh2001.github.io/tree/main/docs/${filePath}`
    },
    footer: {
      message: 'Released under the <a href="/mit">MIT License</a>.',
      copyright: 'Copyright © 2024-present <a href="/resume">Henghua Zhang</a>'
    },
    nav: [
      { text: '首页', link: '/' },
      {
        text: '科研笔记',
        items: [
          { text: 'SDN', items: SDNNoteItems },
        ],
      },
      { text: '其他笔记', items: TechNotes },
      { text: '友情链接', link: '/links' },
      { text: '关于我', link: '/resume' }
    ],
    search: {
      provider: 'local'
    },
    sidebar: {
      '/sdn/': ResearchNotes,
      '/go/': TechNotes,
      '/db/': TechNotes,
      '/interview/': TechNotes,
      '/other/': TechNotes
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/zhh2001' },
      // { icon: 'orcid', link: 'https://orcid.org/0009-0005-9456-8936' },
      { icon: 'qq', link: 'mailto:1652709417@qq.com' },
      { icon: 'leetcode', link: 'https://leetcode.cn/u/zhanghenghua/' },
      { icon: 'csdn', link: 'https://blog.csdn.net/qq_43133192' },
    ]
  },
  sitemap: { hostname: 'https://zhh2001.github.io' }
})
