<script setup>
import Giscus from '@giscus/vue'
import { useData } from 'vitepress'

const { isDark } = useData()
</script>

<template>
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
</template>
