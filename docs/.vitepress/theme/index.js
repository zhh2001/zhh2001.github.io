import DefaultTheme from 'vitepress/theme'
import './index.css'
import './copyright.js'
import Comments from './components/Comments.vue'

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    app.component('Comments', Comments)
  }
}