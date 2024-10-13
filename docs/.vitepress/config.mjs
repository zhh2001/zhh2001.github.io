import { defineConfig } from 'vitepress'

export default defineConfig({
    lang: "zh-Hans",
    title: "张恒华",
    description: "张恒华的个人网站，存放学习笔记与个人简历。研究领域为软件定义网络（SDN，Software Defined Network）。",
    head: [
        ['link', { rel: 'icon', href: '/favicon.ico' }],
        ['meta', { name: 'keywords', content: '张恒华，软件定义网络，网络编程语言，SDN控制器，SDN, P4, Ryu, Mininet' }],
    ],
    themeConfig: {
        returnToTopLabel: '回到顶部',
        sidebarMenuLabel: '菜单',
        darkModeSwitchLabel: '主题',
        darkModeSwitchTitle: '切换到深色模式',
        lightModeSwitchTitle: '切换到浅色模式',
        search: { provider: 'local' },
        outline: { level: 2, label: '页面导航' },
        docFooter: { prev: '上一页', next: '下一页' },

        footer: {
            message: '<a href="https://beian.miit.gov.cn" target="_black">湘ICP备2024042555号</a><br><img src="https://beian.mps.gov.cn/img/logo01.dd7ff50e.png" width="16" height="17" style="vertical-align:middle;display:inline"> <a href="https://beian.mps.gov.cn/#/query/webSearch?code=43112602000261" rel="noreferrer" target="_blank">湘公网安备43112602000261</a>',
            copyright: '基于 <a href="/license">MIT 许可</a> 发布.<br>版权所有 © 2024至今 <a href="/resume">张恒华</a>'
        },

        nav: [
            { text: '首页', link: '/' },
            {
                text: '学习笔记',
                items: [
                    {
                        text: 'SDN',
                        items: [
                            { text: 'Mininet', link: '/sdn/mininet' },
                            { text: 'P4', link: '/sdn/p4' },
                        ]
                    },
                ]
            },
            { text: '简历', link: '/resume' },
        ],

        sidebar: [
            {
                text: 'SDN',
                link: '/sdn/',
                items: [
                    { text: 'Mininet', link: '/sdn/mininet' },
                    { text: 'P4', link: '/sdn/p4' },
                    { text: 'P4 Exercise', link: '/sdn/p4_exercise' },
                ]
            }
        ],

        // socialLinks: [
        //     { icon: 'github', link: 'https://github.com/zhh2001' }
        // ]
    },
    markdown: {
        languageAlias: {
            'p4': 'c++',
            'pseudo': 'c++'
        },
    }
})