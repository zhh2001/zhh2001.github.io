---
layout: page
sidebar: false
aside: false
titleTemplate: 友情链接
---

<script setup>
import {
  VPTeamPage,
  VPTeamPageTitle,
  VPTeamMembers,
  VPTeamPageSection
} from 'vitepress/theme'
import { members, friends } from './.vitepress/theme/links-data.js'
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
</VPTeamPage>

<div style="text-align: center">
  <h2 style="font-size: 24px; font-weight: bold;">申请友链</h2>
  <div style="margin: 16px 0 8px;">
    <strong>申请条件</strong>
    <p>1.请确定贵站可以稳定运营；</p>
    <p>2.博客更新有一定的活跃度；</p>
    <p>3.申请前添加本站至您的友链</p>
  </div>
  <div style="margin-bottom: 8px;">
    <strong>申请须知</strong>
    <p>本站会定期清理长时间无法访问的网站</p>
  </div>
  <div style="margin-bottom: 6px;">
    <strong>申请方式</strong>
    <p>在本页面底部留言，请包含如下信息：</p>
    <p>1.网站名称和地址（必须）</p>
    <p>2.网站描述和Logo（可选）</p>
    <p>3.各社交平台链接（可选）</p>
  </div>
</div>

<Comments />
