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

      <!-- 用户管理（带子菜单） -->
      <el-sub-menu index="user">
        <template #title>
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </template>
        <el-menu-item index="/user/list" @click="handleMenuClick('/user/list')">用户列表</el-menu-item>
        <el-menu-item index="/user/add" @click="handleMenuClick('/user/add')">新增用户</el-menu-item>
      </el-sub-menu>

      <!-- 角色管理 -->
      <el-menu-item index="/role" @click="handleMenuClick('/role')">
        <el-icon><UserFilled /></el-icon>
        <template #title>角色管理</template>
      </el-menu-item>

      <!-- 菜单管理 -->
      <el-menu-item index="/menu" @click="handleMenuClick('/menu')">
        <el-icon><Menu /></el-icon>
        <template #title>菜单管理</template>
      </el-menu-item>

      <!-- 系统设置 -->
      <el-menu-item index="/settings" @click="handleMenuClick('/settings')">
        <el-icon><Setting /></el-icon>
        <template #title>系统设置</template>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
// 引入 Element Plus 图标
import { HomeFilled, User, UserFilled, Menu, Setting, Edit } from '@element-plus/icons-vue'

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
