import { defineAsyncComponent } from 'vue'
import DefaultTheme from 'vitepress/theme'
import './index.css'
import './copyright.js'

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    app.component('Comments', defineAsyncComponent(() => import('./components/Comments.vue')))
  }
}
