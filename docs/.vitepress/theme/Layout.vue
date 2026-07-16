<script setup>
import DefaultTheme from 'vitepress/theme'
import Busuanzi from './components/Busuanzi.vue'

const { Layout } = DefaultTheme
</script>

<template>
  <Layout>
    <template #layout-bottom>
      <Busuanzi />
    </template>
  </Layout>
</template>
