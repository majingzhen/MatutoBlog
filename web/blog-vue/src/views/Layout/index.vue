<template>
  <el-container style="height: 100vh">
    <!-- 侧边栏导航 -->
    <el-aside width="200px" :style="{ background: '#0f172a' }">
      <Sidebar />
    </el-aside>

    <el-container>
      <!-- 顶部导航栏 -->
      <el-header :style="{ background: '#fff', boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }">
        <Topbar />
      </el-header>

      <!-- 主内容区 -->
      <el-main :style="{ padding: '20px', background: '#f1f5f9' }">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
// 引入侧边栏和顶部导航组件
import Sidebar from './Sidebar.vue'
import Topbar from './Topbar.vue'
</script>
