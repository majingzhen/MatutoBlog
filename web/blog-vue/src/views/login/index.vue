<template>
  <div class="login-container">
    <div class="login-wrapper">
      <!-- 左侧品牌区域 -->
      <div class="login-brand">
        <div class="brand-content">
          <h1 class="brand-title">Admin Dashboard</h1>
          <p class="brand-desc">高效、安全的后台管理系统解决方案</p>
          <div class="brand-image">
            <img src="https://picsum.photos/400/300" alt="管理系统展示图" class="illustration" />
          </div>
        </div>
      </div>

      <!-- 右侧登录表单 -->
      <div class="login-form-container">
        <div class="login-card">
          <div class="login-header">
            <h2 class="login-title">账户登录</h2>
            <p class="login-subtitle">请输入您的账号信息</p>
          </div>

          <el-form
              :model="loginForm"
              :rules="loginRules"
              ref="loginFormRef"
              class="login-form"
          >
            <el-form-item prop="username">
              <el-input
                  v-model="loginForm.username"
                  placeholder="请输入用户名"
                  autocomplete="off"
                  size="large"
                  :prefix-icon="User"
              />
            </el-form-item>

            <el-form-item prop="password">
              <el-input
                  v-model="loginForm.password"
                  type="password"
                  placeholder="请输入密码"
                  autocomplete="off"
                  size="large"
                  :prefix-icon="Lock"
                  show-password
              />
            </el-form-item>

            <el-form-item class="form-actions">
              <el-checkbox v-model="loginForm.remember" label="记住我" />
              <el-link type="primary" @click="handleForgotPassword">忘记密码？</el-link>
            </el-form-item>

            <el-form-item>
              <el-button
                  type="primary"
                  style="width: 100%;"
                  @click="handleLogin"
                  :loading="loading"
                  size="large"
              >
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { login } from '@/api/auth' // 引入登录API

// 登录表单数据
const loginForm = reactive({
  username: '',
  password: '',
  remember: false
})

// 表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ]
}

// 表单引用
const loginFormRef = ref(null)
// 加载状态
const loading = ref(false)
// 路由实例
const router = useRouter()

// 页面加载时检查是否有记住的用户名
onMounted(() => {
  const savedUsername = localStorage.getItem('username')
  if (savedUsername) {
    loginForm.username = savedUsername
    loginForm.remember = true
  }
})

// 处理登录
const handleLogin = async () => {
  // 表单验证
  const valid = await loginFormRef.value.validate().catch(err => {
    console.error('表单验证失败:', err)
    return false
  })

  if (!valid) return

  // 登录请求
  loading.value = true
  try {
    // 调用登录API
    const response = await login(loginForm.username, loginForm.password)

    // 保存token
    const { token } = response.data
    localStorage.setItem('token', token)

    // 记住用户名
    if (loginForm.remember) {
      localStorage.setItem('username', loginForm.username)
    } else {
      localStorage.removeItem('username')
    }

    ElMessage.success('登录成功，正在跳转...')
    // 跳转到首页
    setTimeout(() => {
      router.push('/')
    }, 1000)

  } catch (error) {
    ElMessage.error(error.message || '登录失败，请重试')
  } finally {
    loading.value = false
  }
}

// 忘记密码处理
const handleForgotPassword = () => {
  ElMessage.info('请联系管理员重置密码')
}
</script>

<style scoped>
.login-container {
  width: 100%;
  height: 100vh;
  background-color: #f0f2f5;
  overflow: hidden;
}

.login-wrapper {
  display: flex;
  width: 1000px;
  height: 600px;
  margin: 80px auto;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  overflow: hidden;
  background-color: #fff;
}

/* 左侧品牌区域 */
.login-brand {
  flex: 1;
  background: linear-gradient(135deg, #165DFF 0%, #0E42D2 100%);
  color: #fff;
  padding: 40px;
  display: flex;
  align-items: center;
}

.brand-content {
  width: 100%;
}

.brand-title {
  font-size: 32px;
  margin-bottom: 15px;
  font-weight: 700;
}

.brand-desc {
  font-size: 16px;
  opacity: 0.9;
  margin-bottom: 40px;
  line-height: 1.6;
}

.brand-image {
  text-align: center;
}

.illustration {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 右侧登录表单 */
.login-form-container {
  flex: 0 0 400px;
  padding: 60px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.login-card {
  width: 100%;
}

.login-header {
  margin-bottom: 30px;
  text-align: center;
}

.login-title {
  font-size: 24px;
  font-weight: 600;
  color: #1D2129;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: #86909C;
}

.login-form {
  width: 100%;
}

.el-form-item {
  margin-bottom: 20px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 响应式调整 */
@media (max-width: 1024px) {
  .login-wrapper {
    width: 90%;
    flex-direction: column;
    height: auto;
    margin: 40px auto;
  }

  .login-brand {
    flex: none;
    padding: 30px;
  }

  .login-form-container {
    flex: none;
    width: 100%;
    padding: 30px;
  }
}
</style>