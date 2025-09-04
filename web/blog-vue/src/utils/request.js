import axios from 'axios'
import { ElMessage, ElLoading, ElMessageBox } from 'element-plus'
import router from '@/router'

// 创建axios实例
const service = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api', // 从环境变量获取基础URL
    timeout: 10000, // 请求超时时间
    headers: {
        'Content-Type': 'application/json;charset=utf-8'
    }
})

// 加载实例
let loadingInstance = null

// 请求拦截器
service.interceptors.request.use(
    (config) => {
        // 显示加载动画
        loadingInstance = ElLoading.service({
            lock: true,
            text: '加载中...',
            background: 'rgba(0, 0, 0, 0.1)'
        })

        // 获取并添加token
        const token = localStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }

        return config
    },
    (error) => {
        // 关闭加载动画
        if (loadingInstance) {
            loadingInstance.close()
        }

        // 请求错误处理
        ElMessage.error('请求异常，请稍后重试')
        return Promise.reject(error)
    }
)

// 响应拦截器
service.interceptors.response.use(
    (response) => {
        // 关闭加载动画
        if (loadingInstance) {
            loadingInstance.close()
        }

        const res = response.data

        // 根据后端约定的状态码处理
        // 假设成功状态码为200
        if (res.code !== 200) {
            // 错误提示
            ElMessage.error(res.message || '操作失败')

            // 特殊错误码处理
            if (res.code === 401) {
                // 未授权，需要重新登录
                ElMessageBox.confirm(
                    '您的登录已过期，请重新登录',
                    '登录过期',
                    {
                        confirmButtonText: '重新登录',
                        cancelButtonText: '取消',
                        type: 'warning'
                    }
                ).then(() => {
                    // 清除token并跳转到登录页
                    localStorage.removeItem('token')
                    router.push('/login')
                })
            }

            return Promise.reject(new Error(res.message || 'Error'))
        }

        return res
    },
    (error) => {
        // 关闭加载动画
        if (loadingInstance) {
            loadingInstance.close()
        }

        // 网络错误处理
        if (!error.response) {
            ElMessage.error('网络异常，请检查网络连接')
            return Promise.reject(error)
        }

        // HTTP状态码处理
        const status = error.response.status
        switch (status) {
            case 401:
                // 未授权，需要重新登录
                ElMessageBox.confirm(
                    '您的登录已过期，请重新登录',
                    '登录过期',
                    {
                        confirmButtonText: '重新登录',
                        cancelButtonText: '取消',
                        type: 'warning'
                    }
                ).then(() => {
                    localStorage.removeItem('token')
                    router.push('/login')
                })
                break
            case 403:
                ElMessage.error('没有权限执行此操作')
                break
            case 404:
                ElMessage.error('请求的资源不存在')
                break
            case 500:
                ElMessage.error('服务器内部错误')
                break
            default:
                ElMessage.error(`请求错误: ${status}`)
        }

        return Promise.reject(error)
    }
)

export default service
