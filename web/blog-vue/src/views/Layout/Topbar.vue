<template>
  <div class="topbar">
    <!-- 面包屑导航 -->
    <el-breadcrumb separator="/" :style="{ marginLeft: '20px' }">
      <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item v-if="currentPage !== '首页'">{{ currentPage }}</el-breadcrumb-item>
    </el-breadcrumb>

    <!-- 右侧：用户信息和操作 -->
    <div class="topbar-right" :style="{ marginRight: '20px' }">
      <!-- 用户下拉菜单 -->
      <el-dropdown placement="bottom-end">
        <div :style="{ display: 'flex', alignItems: 'center', cursor: 'pointer' }">
          <el-avatar :size="32" :style="{ marginRight: '8px' }">
            <img :src="loginUser?.avater ?? 'https://picsum.photos/200/200'" alt="用户头像" />
          </el-avatar>
          <span>{{ loginUser.username }}</span>
          <el-icon :style="{ marginLeft: '4px' }"><ChevronDown /></el-icon>
        </div>

        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item>
              <el-icon><User /></el-icon>
              <span>个人中心</span>
            </el-dropdown-item>
            <el-dropdown-item>
              <el-icon><Setting /></el-icon>
              <span>账号设置</span>
            </el-dropdown-item>
            <el-dropdown-item divided @click="handleLogout">
              <el-icon><SwitchButton /></el-icon>
              <span>退出登录</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Bell, User, Setting, SwitchButton } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()

const loginUser = computed(() => {
  return JSON.parse(localStorage.getItem('login_user'))
})

// 页面名称映射
const pageNameMap = {
  '/': '首页',
  '/publish': '发布文章',
  '/article': '文章列表',
  '/user': '用户列表',
  '/role': '角色管理',
  '/menu': '菜单管理',
  '/settings': '系统设置'
}

// 当前页面名称
const currentPage = computed(() => {
  return pageNameMap[route.path] || '首页'
})

// 退出登录逻辑
const handleLogout = () => {
  ElMessage.success('退出登录成功')
  localStorage.removeItem('login_user')
  localStorage.removeItem('token')
  router.push('/login')
}
</script>

<style scoped>
.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}
</style>
