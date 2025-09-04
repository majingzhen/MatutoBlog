<template>
  <div class="topbar">
    <!-- 面包屑导航 -->
    <el-breadcrumb separator="/" :style="{ marginLeft: '20px' }">
      <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item>{{ currentPage }}</el-breadcrumb-item>
    </el-breadcrumb>

    <!-- 右侧：用户信息和操作 -->
    <div class="topbar-right" :style="{ marginRight: '20px' }">
      <!-- 通知图标 -->
      <el-icon :style="{ marginRight: '20px', cursor: 'pointer' }">
        <Bell />
      </el-icon>

      <!-- 用户下拉菜单 -->
      <el-dropdown placement="bottom-end">
        <div :style="{ display: 'flex', alignItems: 'center', cursor: 'pointer' }">
          <el-avatar :size="32" :style="{ marginRight: '8px' }">
            <img src="https://picsum.photos/200/200" alt="用户头像" />
          </el-avatar>
          <span>管理员</span>
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
              <el-icon><Logout /></el-icon>
              <span>退出登录</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Bell, User, Setting } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// 当前页面名称
const currentPage = ref('首页')
const router = useRouter()

// 退出登录逻辑
const handleLogout = () => {
  ElMessage.success('退出登录成功')
  // 实际项目中需清除 Token，并重定向到登录页
  // router.push('/login')
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
