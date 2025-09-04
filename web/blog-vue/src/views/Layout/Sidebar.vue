<template>
  <div class="sidebar">
    <!-- 侧边栏标题 -->
    <div class="sidebar-header" :style="{ padding: '20px', color: '#fff', textAlign: 'center', fontSize: '18px' }">
      Vue3 Admin
    </div>

    <!-- 导航菜单 -->
    <el-menu
      default-active="1"
      class="el-menu-vertical-demo"
      :background-color="'#0f172a'"
      :text-color="'#cbd5e1'"
      :active-text-color="'#38bdf8'"
      :collapse-transition="false"
    >
      <!-- 首页 -->
      <el-menu-item index="1">
        <el-icon><HomeFilled /></el-icon>
        <template #title>首页</template>
      </el-menu-item>

      <!-- 发布文章 -->
      <el-menu-item index="2">
        <el-icon><Edit /></el-icon>
        <template #title>发布文章</template>
      </el-menu-item>

      <!-- 用户管理（带子菜单） -->
      <el-sub-menu index="2">
        <template #title>
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </template>
        <el-menu-item index="2-1">用户列表</el-menu-item>
        <el-menu-item index="2-2">新增用户</el-menu-item>
      </el-sub-menu>

      <!-- 角色管理 -->
      <el-menu-item index="3">
        <el-icon><UserFilled /></el-icon>
        <template #title>角色管理</template>
      </el-menu-item>

      <!-- 菜单管理 -->
      <el-menu-item index="4">
        <el-icon><Menu /></el-icon>
        <template #title>菜单管理</template>
      </el-menu-item>

      <!-- 系统设置 -->
      <el-menu-item index="5">
        <el-icon><Setting /></el-icon>
        <template #title>系统设置</template>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
// 引入 Element Plus 图标
import { HomeFilled, User, UserFilled, Menu, Setting } from '@element-plus/icons-vue'
</script>

<style scoped>
.sidebar {
  height: 100%;
  overflow-y: auto;
}
</style>
