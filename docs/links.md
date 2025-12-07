---
layout: page
sidebar: false
aside: false
---

<script setup>
import {
  VPTeamPage,
  VPTeamPageTitle,
  VPTeamMembers,
  VPTeamPageSection
} from 'vitepress/theme'
import Giscus from '@giscus/vue'

import { useData } from 'vitepress'
const { isDark } = useData()

const linkIcon = {
  svg: '<svg t="1759998982284" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="22707" width="64" height="64"><path d="M512 1024a512 512 0 1 1 0-1024 512 512 0 0 1 0 1024z m79.6672-557.8752l-1.979733-1.706667a108.817067 108.817067 0 0 0-11.946667-9.284266l-38.638933 34.542933a62.122667 62.122667 0 0 1 12.629333 8.738133l2.048 1.706667a52.770133 52.770133 0 0 1 0 80.2816L448.238933 674.952533c-24.7808 22.1184-64.853333 22.1184-89.634133 0l-2.048-1.774933a52.701867 52.701867 0 0 1 0-80.2816l47.786667-42.734933a132.642133 132.642133 0 0 1-11.8784-57.344l-73.864534 66.082133a97.621333 97.621333 0 0 0 0 148.343467l1.979734 1.706666c45.533867 40.823467 120.0128 40.823467 165.546666 0l105.540267-94.481066a97.621333 97.621333 0 0 0 0-148.343467z m170.666667-149.2992l-2.048-1.774933c-45.4656-40.823467-120.0128-40.823467-165.546667 0L489.2672 409.6c-45.533867 40.7552-45.533867 107.451733 0 148.2752l1.979733 1.706667c3.822933 3.413333 7.850667 6.417067 11.946667 9.284266l38.638933-34.6112a62.6688 62.6688 0 0 1-12.629333-8.669866l-1.979733-1.706667a52.770133 52.770133 0 0 1 0-80.349867L632.832 349.047467c24.712533-22.186667 64.853333-22.186667 89.565867 0l2.048 1.706666a52.770133 52.770133 0 0 1 0 80.349867l-47.7184 42.734933c8.260267 18.363733 12.219733 37.819733 11.810133 57.344l73.796267-66.082133a97.553067 97.553067 0 0 0 0-148.2752z" fill="#7F7F85" p-id="22708"></path></svg>'
}

const mailIcon = {
  svg: '<svg t="1760000487865" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="15429" width="64" height="64"><path d="M970.56 267.52a85.76 85.76 0 0 0 0-10.24v-2.88l-17.92 1.6 15.68-3.84v-3.52A141.44 141.44 0 0 0 960 227.2a128 128 0 0 0-11.2-20.48 123.84 123.84 0 0 0-103.04-54.72H176.96a123.84 123.84 0 0 0-102.72 55.36 128 128 0 0 0-11.84 22.72 124.8 124.8 0 0 0-6.72 22.72l16 3.2h-16.64a81.92 81.92 0 0 0-1.6 11.84v512a124.48 124.48 0 0 0 124.16 124.16h668.8a124.48 124.48 0 0 0 124.16-124.16v-512zM883.2 304v448a64 64 0 0 1-64 64H205.12a64 64 0 0 1-64-64v-448a64 64 0 0 1 1.92-16l334.4 267.84 2.56 3.2a48.96 48.96 0 0 0 26.88 11.84h6.08a42.24 42.24 0 0 0 31.04-11.84L880.96 288a62.08 62.08 0 0 1 2.24 16zM512 472.96l-290.88-233.6h581.76z" p-id="15430" fill="#7F7F85"></path></svg>'
}

const members = [
  {
    avatar: 'https://blog.ctdxz.com/avatar.webp',
    name: '吃土的小智',
    desc: '我家还蛮大的.jpg',
    links: [
      { icon: linkIcon, link: 'https://blog.ctdxz.com' },
      { icon: mailIcon, link: 'mailto:grtsinry43@outlook.com' },
      { icon: 'qq', link: 'http://wpa.qq.com/msgrd?v=3&uin=3096484572&site=qq&menu=yes' },
      { icon: 'bilibili', link: 'https://space.bilibili.com/476403407' },
      { icon: 'neteasecloudmusic', link: 'https://space.bilibili.com/476403407' },
      { icon: 'github', link: 'https://github.com/xiaozhi-6' },
    ]
  },
  {
    avatar: 'https://static.miaoer.net/logo/avatar.webp',
    name: '喵二の小博客',
    desc: '缘，妙不可言',
    links: [
      { icon: linkIcon, link: 'https://www.miaoer.net' },
      { icon: 'gmail', link: 'mailto:miaoermua@gmail.com' },
      { icon: 'bilibili', link: 'https://space.bilibili.com/41605906' },
      { icon: 'youtube', link: 'https://www.youtube.com/@miaoerchannel' },
      { icon: 'telegram', link: 'https://t.me/miaoer' },
      { icon: 'steam', link: 'https://steamcommunity.com/id/miaoermua' },
    ]
  },
  {
    avatar: 'https://s2.loli.net/2025/02/19/gX19mThYxyUDaZu.gif',
    name: 'Yesord',
    desc: '你相信光吗~~',
    links: [
      { icon: linkIcon, link: 'https://blog.yesord.top' },
      { icon: 'gmail', link: 'xuruolun666@gmail.com' },
      { icon: 'github', link: 'https://github.com/Yesord' },
      { icon: 'csdn', link: 'https://blog.csdn.net/aoliba_believer' },
      { icon: 'bilibili', link: 'https://space.bilibili.com/1527890846' },
      { icon: 'rss', link: 'hhttps://home.yesord.top/atom.xml' },
    ]
  },
  {
    avatar: 'https://q1.qlogo.cn/g?b=qq&nk=3188468169&s=640',
    name: '墨希MoXiify',
    desc: '做一条，逆流的鱼',
    links: [
      { icon: linkIcon, link: 'https://note.moxiify.cn' },
      { icon: 'qq', link: 'mailto:3188468169@qq.com' },
      { icon: 'github', link: 'https://github.com/MoXiaoXi233' },
      { icon: 'rss', link: 'https://note.moxiify.cn/feed/' },
    ]
  },
  {
    avatar: 'https://blog.feng1026.top/avatar.jpg',
    name: '枫落丰源',
    desc: '和你的日常，就是奇迹',
    links: [
      { icon: linkIcon, link: 'https://blog.feng1026.top' },
      { icon: mailIcon, link: 'mailto:qiufengluoye@feng1026.top' },
      { icon: 'bilibili', link: 'https://space.bilibili.com/1735270180' },
      { icon: 'steam', link: 'https://steamcommunity.com/id/qiufengluoye44/' },
      { icon: 'github', link: 'https://github.com/Qiufengluoyes' },
    ]
  },
  {
    avatar: 'https://hzz.cool/favicon.ico',
    name: '何智政个人博客',
    desc: '择善固执 守正出奇',
    links: [
      { icon: linkIcon, link: 'https://hzz.cool' },
      { icon: 'gmail', link: 'mailto:dexter.ho.cn@gmail.com' },
      { icon: 'github', link: 'https://github.com/hezhizheng' },
      { icon: 'instagram', link: 'https://www.instagram.com/dexter_ho_cn' },
      { icon: 'telegram', link: 'https://t.me/dexterho' },
      { icon: 'sinaweibo', link: 'https://weibo.com/u/5675317400' },
    ]
  },
  {
    avatar: 'https://img.cdn1.vip/i/68bbdee513eb9_1757142757.webp',
    name: 'LYEy_isine个人博客',
    desc: '花海无一日,少年踏自来',
    links: [
      { icon: linkIcon, link: 'https://caiyifeng.top' },
      { icon: 'qq', link: 'https://qm.qq.com/q/Uz2Vg3uzC2' },
      { icon: 'bilibili', link: 'https://space.bilibili.com/1358278810' },
      { icon: 'steam', link: 'https://steamcommunity.com/id/LYEy_isine/' },
      { icon: 'github', link: 'https://github.com/caiyifeng0705' },
      { icon: 'rss', link: 'https://caiyifeng.top/rss.xml' },
    ]
  },
  {
    avatar: 'https://img520.com/iNvl5J.png',
    name: '云野阁',
    desc: '闲云野鹤，八方逍遥',
    links: [
      { icon: linkIcon, link: 'https://yyg.js.cool/' },
      { icon: 'github', link: 'https://github.com/MarsperL' },
      { icon: 'rss', link: 'https://yyg.js.cool/rss2.xml' },
    ]
  },
  {
    avatar: 'https://www.plsshenyun.top/ico/ico.png',
    name: '望夜菌的小窝',
    desc: 'change and challenge!',
    links: [
      { icon: linkIcon, link: 'https://www.plsshenyun.top/' },
      { icon: mailIcon, link: 'mailto:blog@plsshenyun.top' },
      { icon: 'rss', link: 'https://www.plsshenyun.top/feed' },
    ]
  },
  {
    avatar: 'https://www.timochan.cn/api/objects/icon/9s6tbcvax674yv2m88.jpg',
    name: 'TimochanのBlog',
    desc: 'Let\'s start learning',
    links: [
      { icon: linkIcon, link: 'https://www.timochan.cn' },
      { icon: mailIcon, link: 'mailto:i@timochan.cn' },
      { icon: 'github', link: 'https://github.com/ttimochan' },
      { icon: 'x', link: 'https://twitter.com/RefRebel' },
      { icon: 'rss', link: 'https://www.timochan.cn/feed' },
    ]
  },
  {
    avatar: 'https://image.kong.college/i/2025/09/09/spq9n6.png',
    name: '28.7',
    desc: '空山不见人，但闻人语响',
    links: [
      { icon: linkIcon, link: 'https://blog.kong.college' },
    ]
  },
  {
    avatar: 'https://mccsjs.eu.org/img/head.jpg',
    name: 'mccsjs',
    desc: '点一盏灯，等一个迷路的夜',
    links: [
      { icon: linkIcon, link: 'https://mccsjs.eu.org' },
      { icon: 'qq', link: 'https://res.abeim.cn/api/qq/?qq=3505591664' },
      { icon: 'wechat', link: 'https://mccsjs.eu.org/img/wx.jpg' },
      { icon: 'bilibili', link: 'https://space.bilibili.com/209190096' },
      { icon: 'github', link: 'https://github.com/mccsjs' },
      { icon: 'rss', link: 'https://mccsjs.eu.org/rss2.xml' },
    ]
  },
  {
    avatar: 'https://cdn.tulan.cyou/noah/2025/07/26/b_5255bd8e31082e2a4332c4eb4172686e.jpg',
    name: 'NoahのBlog',
    desc: '努力去发光，而不是被照亮',
    links: [
      { icon: linkIcon, link: 'https://blog.noah0932.top' },
      { icon: mailIcon, link: 'mailto:noah@noah0932.top' },
      { icon: 'github', link: 'https://github.com/Noah0932' },
      { icon: 'rss', link: 'https://blog.noah0932.top/rss.xml' },
    ]
  },
  {
    avatar: 'https://importmengjie.github.io/images/avatar.jpg',
    name: 'ImportMengjie',
    desc: '一个车端C++软件开发程序猿',
    links: [
      { icon: linkIcon, link: 'https://importmengjie.github.io' },
      { icon: mailIcon, link: 'mailto:limengjie@hotmail.com' },
      { icon: 'github', link: 'https://github.com/ImportMengjie' },
      { icon: 'rss', link: 'https://importmengjie.github.io/index.xml' },
    ]
  },
  {
    avatar: 'https://free.picui.cn/free/2025/09/30/68dbd5869ed39.jpg',
    name: '裕裕裕的小破宅',
    desc: '一个充满青春活力的技术博客',
    links: [
      { icon: linkIcon, link: 'https://yu-blog.top' },
      { icon: 'gmail', link: 'mailto:https://github.com/rossiniortensia-ops' },
      { icon: 'sinaweibo', link: 'https://weibo.com/u/7995211964' },
      { icon: 'github', link: 'https://github.com/rossiniortensia-ops' },
    ]
  },
  {
    avatar: 'https://i1-cdn.feizhuqwq.com/img-assets/logo/faviconHD.jpg',
    name: 'feizhuqwq',
    desc: '因为不可能，所以才值得相信',
    links: [
      { icon: linkIcon, link: 'https://blog.feizhuqwq.com' },
      { icon: mailIcon, link: 'mailto:me@feizhuqwq.com' },
      { icon: 'rss', link: 'https://blog.feizhuqwq.com/feed/' },
    ]
  },
  {
    avatar: 'https://blog.tovoao.cn/favicon.ico',
    name: '鹤归博客',
    desc: '总有人间一两风<br>填我十万八千梦',
    links: [
      { icon: linkIcon, link: 'https://blog.tovoao.cn' },
      { icon: mailIcon, link: 'mailto:981739185@qq.com' },
      { icon: 'qq', link: 'https://wpa.qq.com/msgrd?v=3&uin=981739185&site=qq&menu=yes' },
    ]
  },
  {
    avatar: 'https://typonotes.com/logo/avatar.png',
    name: '老麦的书房',
    desc: 'Golang、 云原生、 DevOps、 可视化追踪',
    links: [
      { icon: linkIcon, link: 'https://typonotes.com' },
      { icon: 'github', link: 'https://github.com/tangx' },
      { icon: 'x', link: 'https://twitter.com/tangx' },
      { icon: 'rss', link: 'https://typonotes.com/index.xml' },
    ]
  },
  {
    avatar: 'https://chenmingyong.cn/static/4aa6f9aeeee31495a9fb2cc6d2f7a1ca.jpg',
    name: '陈明勇的博客',
    desc: '一名热爱技术、乐于分享的开发者，同时也是开源爱好者',
    links: [
      { icon: linkIcon, link: 'https://chenmingyong.cn' },
      { icon: 'gmail', link: 'mailto:chenmingyong1999@gmail.com' },
      { icon: 'github', link: 'https://github.com/chenmingyong0423' },
      { icon: 'wechat', link: 'https://chenmingyong.cn/static/wx-qrcode.jpg' },
      { icon: 'zhihu', link: 'https://www.zhihu.com/people/chenmingyong-code' },
      { icon: 'juejin', link: 'https://juejin.cn/user/4174180683088269' },
    ]
  },
  {
    avatar: 'https://pic1.imgdb.cn/item/68b512ef58cb8da5c8689e8c.webp',
    name: '一世繁华',
    title: '水下机器人领域从业者',
    desc: '分享一二',
    links: [
      { icon: linkIcon, link: 'https://blog.hantaotao.top' },
      { icon: mailIcon, link: 'mailto:tohantao@outlook.com' },
    ]
  },
  {
    avatar: 'https://bu.dusays.com/2024/10/25/671b2438203a6.gif',
    name: 'Elykia',
    org: 'ZZU',
    orgLink: 'https://www.zzu.edu.cn/',
    desc: '致以无暇之人',
    links: [
      { icon: linkIcon, link: 'https://blog.elykia.cn' },
      { icon: 'qq', link: 'mailto:elykia@qq.com' },
      { icon: 'bilibili', link: 'https://space.bilibili.com/1451399239' },
      { icon: 'gitee', link: 'https://gitee.com/Elykia' },
      { icon: 'github', link: 'https://github.com/elykia-cn' },
      { icon: 'telegram', link: 'https://t.me/Elykia_cn' },
    ]
  },
  {
    avatar: 'https://blog.grtsinry43.com/favicon.ico',
    name: 'Grtsinry43\'s Blog',
    title: '',
    org: 'CSU',
    orgLink: 'https://www.csu.edu.cn',
    desc: '总之岁月漫长，然而值得等待',
    links: [
      { icon: linkIcon, link: 'https://blog.grtsinry43.com' },
      { icon: mailIcon, link: 'mailto:grtsinry43@outlook.com' },
      { icon: 'github', link: 'https://github.com/grtsinry43' },
      { icon: 'rss', link: 'https://blog.grtsinry43.com/feed' },
    ]
  },
  {
    avatar: 'https://youpai.roccoshi.top/avatar.jpg',
    name: 'Moreality\'s Blog',
    title: 'Developer',
    org: 'Amazon',
    orgLink: 'https://www.amazon.com',
    desc: 'The singularity is nearer.',
    links: [
      { icon: linkIcon, link: 'https://moreality.net' },
      { icon: 'gmail', link: 'mailto:imroccoshi@gmail.com' },
      { icon: 'github', link: 'https://github.com/Lincest' },
      { icon: 'bilibili', link: 'https://space.bilibili.com/15255409' },
      { icon: 'x', link: 'https://x.com/himoreality' },
      { icon: 'rss', link: 'https://moreality.net/atom.xml' },
    ]
  },
  {
    avatar: 'https://res.strikefreedom.top/static_res/blog/figures/avatar.png',
    name: 'Strike Freedom',
    org: 'GNET',
    orgLink: 'https://gnet.host',
    desc: '潘少的博客、主页、技术分享',
    links: [
      { icon: linkIcon, link: 'https://strikefreedom.top' },
      { icon: 'github', link: 'https://github.com/panjf2000' },
      { icon: 'zhihu', link: 'https://www.zhihu.com/people/andy_pan' },
      { icon: 'discord', link: 'https://discord.gg/UyKD7NZcfH' },
      { icon: 'instagram', link: 'https://instagram.com/panjf2000' },
      { icon: 'mastodon', link: 'https://mastodon.social/@andypan' },
    ]
  },
]

const friends = [
  {
    avatar: 'https://xiongyujie.cn/img/xiongyujie.jpg',
    name: '熊玉洁',
    title: '副教授/硕士生导师',
    org: 'SUES',
    orgLink: 'https://www.sues.edu.cn/',
    desc: '华东师范大学博士',
    links: [
      { icon: linkIcon, link: 'https://xiongyujie.cn/' },
      { icon: mailIcon, link: 'mailto:xiong@sues.edu.cn' },
      { icon: 'github', link: 'https://github.com/X-Lab-CN' },
    ]
  },
  {
    avatar: 'https://seee.sues.edu.cn/_upload/article/images/a2/b1/5596c5f54495bd320eea8ba7dbd9/70e432c7-4d57-4535-9cbb-cdbd2e4741da.jpg',
    name: '陈珏',
    title: '讲师/硕士生导师',
    org: 'SUES',
    orgLink: 'https://www.sues.edu.cn/',
    desc: '华东师范大学博士',
    links: [
      { icon: linkIcon, link: 'https://seee.sues.edu.cn/fb/d9/c20786a195545/page.htm' },
      { icon: mailIcon, link: 'mailto:jadeschen@sues.edu.cn' },
    ]
  },
  {
    avatar: 'https://vovyh.github.io/images/avatar.jpg',
    name: '邬雨航',
    title: '硕士生',
    org: 'SUES',
    orgLink: 'https://www.sues.edu.cn/',
    desc: '大模型爱好者',
    links: [
      { icon: linkIcon, link: 'https://vovyh.github.io/' },
      { icon: mailIcon, link: 'mailto:m325124620@sues.edu.cn' },
      { icon: 'gmail', link: 'mailto:vovyh0514@gmail.com' },
      { icon: 'github', link: 'https://github.com/VovyH' },
      { icon: 'leetcode', link: 'https://leetcode.cn/u/wuyuhangwinner/' },
      { icon: 'bilibili', link: 'https://space.bilibili.com/1054622831' },
    ]
  },
  {
    avatar: 'https://brandpeng.github.io/static/assets/img/photo.png',
    name: '彭海东',
    title: '硕士生',
    org: 'SUES',
    orgLink: 'https://www.sues.edu.cn/',
    desc: '入侵检测领域大佬',
    links: [
      { icon: linkIcon, link: 'https://brandpeng.github.io/' },
      { icon: mailIcon, link: 'mailto:m325124607@sues.edu.cn' },
      { icon: 'github', link: 'https://github.com/BrandPeng' },
    ]
  },
]

const templates = [
  {
    avatar: 'https://zhh2001.github.io/avatar.jpg',
    name: '[网站名称]',
    title: '[角色]',
    org: '[所在组织]',
    orgLink: '',
    desc: '[网站描述]',
    links: [
      { icon: 'gmail', ariaLabel: 'Gmail' },
      { icon: 'qq', ariaLabel: 'QQ' },
      { icon: 'wechat', ariaLabel: 'WeChat' },
      { icon: 'sinaweibo', ariaLabel: '新浪微博' },
      { icon: 'zhihu', ariaLabel: '知乎' },
      { icon: 'bilibili', ariaLabel: 'B站' },
      { icon: 'tiktok', ariaLabel: 'TikTok' },
      { icon: 'youtube', ariaLabel: 'YouTube' },
      { icon: 'facebook', ariaLabel: 'Facebook' },
    ]
  }
]
</script>

<VPTeamPage>
  <VPTeamPageTitle>
    <template #title>友情链接</template>
    <template #lead>各路大佬</template>
  </VPTeamPageTitle>
  <VPTeamMembers size="small" :members="members" />
  <VPTeamPageSection>
    <template #title>成长路上</template>
    <template #lead>朋友、同学、老师</template>
    <template #members>
      <VPTeamMembers size="small" :members="friends" />
    </template>
  </VPTeamPageSection>
  <VPTeamPageSection>
    <template #title>申请友链</template>
    <template #lead>
      <div style="margin-bottom: 8px;">
        <strong>申请条件</strong>
        <p>1.请确定贵站可以稳定运营；</p>
        <p>2.博客更新有一定的活跃度；</p>
        <p>3.申请前先添加本站至您的友链</p>
      </div>
      <div style="margin-bottom: 8px;">
        <strong>申请须知</strong>
        <p>本站会定期清理长时间无法访问的网站</p>
      </div>
      <div style="margin-bottom: 6px;">
        <strong>申请方式</strong>
        <p>可通过 <a href="https://github.com/zhh2001/zhh2001.github.io/issues" target="_blank">Issue</a> 或者邮件联系我，请包含如下信息：</p>
        <p>1.网站名称和地址（必须）</p>
        <p>2.网站描述和Logo（可选）</p>
        <p>3.各个社交平台及链接（可选）</p>
        <p>4.所在组织及角色（可选）</p>
      </div>
      <div style="margin-bottom: 8px;">
        <strong>这些信息将像下面这样展示</strong>
      </div>
    </template>
    <template #members>
      <VPTeamMembers size="medium" :members="templates" />
    </template>
  </VPTeamPageSection>
</VPTeamPage>

<Giscus
  id="comments"
  repo="zhh2001/zhh2001.github.io"
  repoId="R_kgDOK7hmeA"
  category="Announcements"
  categoryId="DIC_kwDOK7hmeM4CobaL"
  mapping="pathname"
  strict="1"
  reactionsEnabled="1"
  emitMetadata="0"
  inputPosition="top"
  :theme="isDark ? 'dark' : 'light'"
  lang="en"
/>
