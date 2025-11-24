import { defineConfig } from 'vitepress'

const siteMapIcon = {
  svg: '<svg t="1760360619836" class="icon" viewBox="0 0 1280 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="3431" width="64" height="64"><path d="M256 704H64c-35.34 0-64 28.66-64 64v192c0 35.34 28.66 64 64 64h192c35.34 0 64-28.66 64-64v-192c0-35.34-28.66-64-64-64z m-48-160h384v96h96v-96h384v96h96v-115.18c0-42.34-34.46-76.82-76.82-76.82H688v-128h80c35.34 0 64-28.66 64-64V64c0-35.34-28.66-64-64-64H512c-35.34 0-64 28.66-64 64v192c0 35.34 28.66 64 64 64h80v128H188.82C146.46 448 112 482.46 112 524.82V640h96v-96z m528 160h-192c-35.34 0-64 28.66-64 64v192c0 35.34 28.66 64 64 64h192c35.34 0 64-28.66 64-64v-192c0-35.34-28.66-64-64-64z m480 0h-192c-35.34 0-64 28.66-64 64v192c0 35.34 28.66 64 64 64h192c35.34 0 64-28.66 64-64v-192c0-35.34-28.66-64-64-64z" p-id="3432" fill="#7F7F85"></path></svg>'
}

const SDNNoteItems = [
  { text: 'P4 Language', link: '/sdn/p4' },
  { text: 'P4 Exercise', link: '/sdn/p4_exercise' },
  { text: 'P4 INT', link: '/sdn/int' },
  { text: 'P4 Utils', link: '/sdn/p4_utils' },
  { text: 'Mininet', link: '/sdn/mininet' },
  { text: 'iPerf', link: '/sdn/iperf' },
  {
    text: 'Writing', link: '/sdn/writing/', items: [
      { text: 'Formula', link: '/sdn/writing/formula' },
      { text: 'Algorithm', link: '/sdn/writing/algorithm' },
      { text: 'Table', link: '/sdn/writing/table' },
      { text: 'Bibliography', link: '/sdn/writing/bibliography' }
    ]
  },
  { text: 'CCF', link: '/sdn/ccf' }
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
  '笔记', '学习笔记', '科研笔记', '读书笔记', '阅读笔记',
  '文献阅读', '论文阅读',
  '软件定义网络', '可编程数据平面', '带内网络遥测',
  'Go语言', 'Go协程', 'Go开发',
  'SDN控制器',
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
  'Ryu',
  'Go', 'GoLang',
  'RPC', 'Remote Procedure Call', 'gRPC',
  'Goroutine', 'GMP',
  'Gin',
  'Docker',
  'Eino',
  'Database',
  'Redis', 'Remote Dictionary Server',
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
  ignoreDeadLinks: true,
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
    container: {
      dangerLabel: '危险',
      detailsLabel: '详细信息',
      infoLabel: '信息',
      tipLabel: '提示',
      warningLabel: '警告'
    },
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
    outlineTitle: '页面导航',
    darkModeSwitchLabel: '主题',
    lightModeSwitchTitle: '切换到浅色模式',
    darkModeSwitchTitle: '切换到深色模式',
    sidebarMenuLabel: '菜单',
    returnToTopLabel: '回到顶部',
    docFooter: { prev: '上一页', next: '下一页' },
    editLink: {
      pattern: ({ filePath }) => `https://github.com/zhh2001/notes/tree/main/docs/${filePath}`,
      text: '在 GitHub 上编辑此页面'
    },
    footer: {
      message: '基于 <a href="/mit">MIT 许可</a> 发布',
      copyright: '版权所有 © 2024至今 <a href="/resume">张恒华</a>'
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
      provider: 'local',
      options: {
        translations: {
          button: { buttonText: '搜索文档' },
          modal: {
            displayDetails: '显示详细内容',
            resetButtonTitle: '清除查询条件',
            noResultsText: '无法找到相关结果',
            footer: {
              selectText: '选择',
              navigateText: '切换',
              closeText: '关闭'
            }
          }
        }
      }
    },
    sidebar: {
      '/sdn/': ResearchNotes,
      '/go/': TechNotes,
      '/db/': TechNotes,
      '/interview/': TechNotes,
      '/other/': TechNotes
    },
    notFound: {
      code: '404',
      linkText: '返回首页',
      quote: '请检查页面路径是否正确',
      title: '页面未找到'
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/zhh2001' },
      // { icon: 'orcid', link: 'https://orcid.org/0009-0005-9456-8936' },
      { icon: 'qq', link: 'mailto:1652709417@qq.com' },
      // { icon: 'leetcode', link: 'https://leetcode.cn/u/zhanghenghua/' },
      { icon: 'csdn', link: 'https://blog.csdn.net/qq_43133192' },
    ]
  },
  sitemap: { hostname: 'https://zhh2001.github.io' }
})
