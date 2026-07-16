import { defineAsyncComponent } from 'vue'
import DefaultTheme from 'vitepress/theme'
import Layout from './Layout.vue'
import './index.css'
import './copyright.js'

export default {
  extends: DefaultTheme,
  Layout,
  enhanceApp({ app }) {
    app.component('Comments', defineAsyncComponent(() => import('./components/Comments.vue')))
  }
}
