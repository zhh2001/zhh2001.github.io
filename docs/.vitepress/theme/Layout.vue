<script setup>
import DefaultTheme from 'vitepress/theme'
import SiteUptime from './components/SiteUptime.vue'

const { Layout } = DefaultTheme
</script>

<template>
  <Layout>
    <template #layout-bottom>
      <SiteUptime />
    </template>
  </Layout>
</template>
