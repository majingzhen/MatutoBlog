<template>
  <div class="sidebar">
    <!-- 侧边栏标题 -->
    <div class="sidebar-header" :style="{ padding: '20px', color: '#fff', textAlign: 'center', fontSize: '18px' }">
      Vue3 Admin
    </div>

    <!-- 导航菜单 -->
    <el-menu
      :default-active="currentRoute"
      class="el-menu-vertical-demo"
      :background-color="'#0f172a'"
      :text-color="'#cbd5e1'"
      :active-text-color="'#38bdf8'"
      :collapse-transition="false"
      router
    >
      <!-- 首页 -->
      <el-menu-item index="/" @click="handleMenuClick('/')">
        <el-icon><HomeFilled /></el-icon>
        <template #title>首页</template>
      </el-menu-item>

      <!-- 发布文章 -->
      <el-menu-item index="/publish" @click="handleMenuClick('/publish')">
        <el-icon><Edit /></el-icon>
        <template #title>发布文章</template>
      </el-menu-item>

      <!-- 文章管理（带子菜单） -->
      <el-sub-menu index="article">
        <template #title>
          <el-icon><Document /></el-icon>
          <span>文章管理</span>
        </template>
        <el-menu-item index="/article" @click="handleMenuClick('/article')">文章列表</el-menu-item>
        <el-menu-item index="/publish" @click="handleMenuClick('/publish')">发布文章</el-menu-item>
      </el-sub-menu>
      <el-menu-item index="/attach" @click="handleMenuClick('/attach')">
        <el-icon><Document /></el-icon>
        <template #title>附件管理</template>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
// 引入 Element Plus 图标
import { HomeFilled, User, UserFilled, Menu, Setting, Edit, Document } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

// 当前活动路由
const currentRoute = computed(() => {
  return route.path
})

// 处理菜单点击
const handleMenuClick = (path) => {
  console.log('点击菜单:', path)
  if (route.path !== path) {
    router.push(path).catch(err => {
      console.log('路由跳转:', err)
    })
  }
}

onMounted(() => {
  console.log('当前路由:', route.path)
})
</script>

<style scoped>
.sidebar {
  height: 100%;
  overflow-y: auto;
}
</style>
