<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <div class="stats" style="display: flex; gap: 20px; margin-bottom: 20px;">
      <el-card style="flex: 1;">
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <p style="color: #666; font-size: 14px;">总用户数</p>
            <h3 style="font-size: 24px; margin-top: 5px;">{{ userCount }}</h3>
          </div>
          <el-icon style="font-size: 36px; color: #409eff;">
            <UserFilled />
          </el-icon>
        </div>
      </el-card>

      <el-card style="flex: 1;">
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <p style="color: #666; font-size: 14px;">今日新增</p>
            <h3 style="font-size: 24px; margin-top: 5px;">{{ todayNew }}</h3>
          </div>
          <el-icon style="font-size: 36px; color: #27ae60;">
            <Plus />
          </el-icon>
        </div>
      </el-card>

      <el-card style="flex: 1;">
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <p style="color: #666; font-size: 14px;">活跃用户</p>
            <h3 style="font-size: 24px; margin-top: 5px;">{{ activeUsers }}</h3>
          </div>
          <el-icon style="font-size: 36px; color: #f39c12;">
            <Activity />
          </el-icon>
        </div>
      </el-card>

      <el-card style="flex: 1;">
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <p style="color: #666; font-size: 14px;">系统状态</p>
            <h3 style="font-size: 24px; margin-top: 5px; color: #27ae60;">正常</h3>
          </div>
          <el-icon style="font-size: 36px; color: #16a085;">
            <CheckCircle />
          </el-icon>
        </div>
      </el-card>
    </div>

    <!-- 最近操作记录 -->
    <el-card style="margin-bottom: 20px;">
      <div slot="header" style="display: flex; justify-content: space-between; align-items: center;">
        <h3>最近操作记录</h3>
        <el-button type="text">查看全部</el-button>
      </div>
      <el-table :data="recentLogs" border>
        <el-table-column label="ID" prop="id" width="80" />
        <el-table-column label="操作人" prop="operator" />
        <el-table-column label="操作内容" prop="content" />
        <el-table-column label="操作时间" prop="time" width="180" />
        <el-table-column label="IP地址" prop="ip" width="140" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { UserFilled, Plus } from '@element-plus/icons-vue'

// 统计数据
const userCount = ref(245)
const todayNew = ref(12)
const activeUsers = ref(189)

// 最近操作记录
const recentLogs = ref([])

// 页面加载时获取数据
onMounted(() => {
  // 模拟获取最近操作记录
  recentLogs.value = [
    {
      id: 1,
      operator: 'admin',
      content: '新增用户 user123',
      time: '2023-09-04 09:23:45',
      ip: '192.168.1.100'
    },
    {
      id: 2,
      operator: 'user2',
      content: '修改角色权限',
      time: '2023-09-04 08:15:32',
      ip: '192.168.1.105'
    },
    {
      id: 3,
      operator: 'admin',
      content: '删除过期数据',
      time: '2023-09-03 16:45:12',
      ip: '192.168.1.100'
    }
  ]
})
</script>
