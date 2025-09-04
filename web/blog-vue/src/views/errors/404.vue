<template>
  <div class="not-found-container">
    <div class="not-found-content">
      <!-- 404 数字图标 -->
      <div class="error-number">
        <span class="digit">4</span>
        <span class="circle">
          <div class="eye"></div>
          <div class="eye"></div>
        </span>
        <span class="digit">4</span>
      </div>
      
      <!-- 错误信息 -->
      <h1 class="error-title">页面不存在</h1>
      <p class="error-message">
        抱歉，您访问的页面不存在或已被移除
      </p>
      
      <!-- 操作按钮 -->
      <div class="error-actions">
        <el-button 
          type="primary" 
          size="large"
          @click="handleGoHome"
        >
          <el-icon><Home /></el-icon>
          返回首页
        </el-button>
        <el-button 
          size="large"
          @click="handleGoBack"
          :disabled="!canGoBack"
        >
          <el-icon><ArrowLeft /></el-icon>
          返回上一页
        </el-button>
      </div>
    </div>
    
    <!-- 装饰元素 -->
    <div class="decor-element"></div>
    <div class="decor-element"></div>
    <div class="decor-element"></div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Home, ArrowLeft } from '@element-plus/icons-vue'

const router = useRouter()
const canGoBack = ref(false)

// 检查是否可以返回上一页
onMounted(() => {
  canGoBack.value = window.history.length > 1
})

// 返回首页
const handleGoHome = () => {
  router.push('/')
}

// 返回上一页
const handleGoBack = () => {
  window.history.back()
}
</script>

<style scoped>
.not-found-container {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #f9fafb;
  position: relative;
  overflow: hidden;
}

.not-found-content {
  text-align: center;
  position: relative;
  z-index: 10;
  max-width: 600px;
  padding: 0 20px;
}

.error-number {
  font-size: 120px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 30px;
  color: #165DFF;
}

.digit {
  position: relative;
  text-shadow: 0 5px 10px rgba(22, 93, 255, 0.1);
}

.circle {
  width: 100px;
  height: 100px;
  border: 20px solid #165DFF;
  border-radius: 50%;
  margin: 0 15px;
  position: relative;
  box-shadow: 0 5px 15px rgba(22, 93, 255, 0.1);
  animation: float 3s ease-in-out infinite;
}

.eye {
  width: 20px;
  height: 20px;
  background-color: #165DFF;
  border-radius: 50%;
  position: absolute;
  top: 30px;
}

.eye:first-child {
  left: 25px;
}

.eye:last-child {
  right: 25px;
}

.error-title {
  font-size: 32px;
  font-weight: 600;
  color: #1D2129;
  margin-bottom: 15px;
}

.error-message {
  font-size: 16px;
  color: #86909C;
  margin-bottom: 40px;
  line-height: 1.6;
}

.error-actions {
  display: flex;
  gap: 15px;
  justify-content: center;
}

/* 装饰元素 */
.decor-element {
  position: absolute;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(22, 93, 255, 0.1) 0%, rgba(255, 255, 255, 0) 70%);
}

.decor-element:nth-child(1) {
  width: 400px;
  height: 400px;
  top: -100px;
  left: -100px;
}

.decor-element:nth-child(2) {
  width: 300px;
  height: 300px;
  bottom: -50px;
  right: 100px;
}

.decor-element:nth-child(3) {
  width: 200px;
  height: 200px;
  top: 50%;
  right: 50px;
}

/* 动画效果 */
@keyframes float {
  0% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-15px);
  }
  100% {
    transform: translateY(0px);
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .error-number {
    font-size: 80px;
  }
  
  .circle {
    width: 70px;
    height: 70px;
    border-width: 15px;
  }
  
  .eye {
    width: 15px;
    height: 15px;
    top: 20px;
  }
  
  .error-title {
    font-size: 24px;
  }
  
  .error-actions {
    flex-direction: column;
    align-items: center;
  }
}
