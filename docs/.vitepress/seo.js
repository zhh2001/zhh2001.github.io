export const SITE_URL = 'https://zhh2001.github.io'

const PAGE_META = new Set([
  'og:title', 'og:description', 'og:url', 'og:type',
  'twitter:title', 'twitter:description'
])

export function transformPageData(pageData, { siteConfig }) {
  const head = pageData.frontmatter.head || []
  // 构建时全局 head 不写入客户端数据，保留默认值以便离开 SDN 页面时恢复。
  let defaults = (siteConfig.site.head || []).filter(([tag, attrs]) =>
    tag === 'meta' && PAGE_META.has(attrs.property || attrs.name)
  )

  if (pageData.relativePath.startsWith('sdn/')) {
    const path = pageData.relativePath
      .replace(/(^|\/)index\.md$/, '$1')
      .replace(/\.md$/, siteConfig.cleanUrls ? '' : '.html')
    const url = new URL(path, new URL(siteConfig.site.base, SITE_URL)).href
    const title = pageData.title || siteConfig.site.title
    const description = pageData.description || siteConfig.site.description
    defaults = [
      ['link', { rel: 'canonical', href: url }],
      ['meta', { property: 'og:title', content: title }],
      ['meta', { property: 'og:description', content: description }],
      ['meta', { property: 'og:url', content: url }],
      ['meta', { property: 'og:type', content: pageData.relativePath.endsWith('/index.md') ? 'website' : 'article' }],
      ['meta', { name: 'twitter:title', content: title }],
      ['meta', { name: 'twitter:description', content: description }]
    ]
  }

  // 页面 head 中的显式配置优先于默认值。
  const missing = defaults.filter(([tag, attrs]) => !head.some(([existingTag, existingAttrs]) => {
    if (tag !== existingTag) return false
    if (tag === 'link') return existingAttrs.rel === attrs.rel
    return attrs.property
      ? existingAttrs.property === attrs.property
      : existingAttrs.name === attrs.name
  }))

  // 写入页面数据，让首次加载与站内导航使用相同的元数据。
  pageData.frontmatter.head = [...missing, ...head]
}
