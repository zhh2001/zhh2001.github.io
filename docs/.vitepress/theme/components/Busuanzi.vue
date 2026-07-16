<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  const script = document.createElement('script')
  script.async = true
  script.src = '//busuanzi.ibruce.info/busuanzi/2.3/busuanzi.pure.mini.js'
  document.head.appendChild(script)
})
</script>

<template>
  <div class="busuanzi-counter">
    <span id="busuanzi_container_site_pv" style="display:none;">
      <span id="busuanzi_value_site_pv"></span> visits
    </span>
    <span id="busuanzi_container_site_uv" style="display:none;">
      <span class="sep">|</span> <span id="busuanzi_value_site_uv"></span> visitors
    </span>
  </div>
</template>

<style scoped>
.busuanzi-counter {
  text-align: center;
  padding: 0 0 24px;
  font-size: 0.85rem;
  color: var(--vp-c-text-2);
}

.sep {
  margin: 0 8px;
  color: var(--vp-c-divider);
}
</style>
