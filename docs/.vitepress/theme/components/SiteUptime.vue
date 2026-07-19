<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vitepress'

const launched = new Date('2023-12-19')
const years = ref(0)
const restDays = ref(0)

function update() {
  const now = new Date()
  const diffMs = now.getTime() - launched.getTime()
  const totalDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  years.value = Math.floor(totalDays / 365)
  restDays.value = totalDays % 365
}

update()

const router = useRouter()

function loadBusuanzi() {
  delete window.busuanziRequestSent
  document.querySelectorAll('script[src*="busuanzi"]').forEach(s => s.remove())
  const script = document.createElement('script')
  script.async = true
  script.src = '//cdn.busuanzi.cc/busuanzi/3.6.9/busuanzi.min.js'
  document.head.appendChild(script)
}

onMounted(() => {
  loadBusuanzi()
  router.onAfterRouteChange = loadBusuanzi
})

onUnmounted(() => {
  router.onAfterRouteChange = null
})
</script>

<template>
  <div class="footer-meta">
    <span class="uptime">{{ years }} yr {{ restDays }} d online</span>
    <span class="sep">|</span>
    <span id="busuanzi_site_pv">...</span> visits
  </div>
</template>

<style scoped>
.footer-meta {
  text-align: center;
  padding: 0 0 20px;
  font-size: 0.875rem;
  color: var(--vp-c-text-2);
}

.sep {
  margin: 0 10px;
  color: var(--vp-c-text-3);
}
</style>
