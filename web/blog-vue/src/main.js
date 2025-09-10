// src/main.js
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import 'element-plus/dist/index.css'

// 引入全局样式
import './styles/global.css'

// 引入 Element Plus 图标库
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 引入 Iconify Vue 组件
import { Icon } from '@iconify/vue'

const app = createApp(App)

// 全局注册所有 Element Plus 图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

// 全局注册 Iconify 图标组件
app.component('Icon', Icon)

app.use(createPinia())
app.use(router)
app.mount('#app')