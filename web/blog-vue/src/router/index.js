import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/views/Layout/index.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // 登录页路由
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/index.vue'),
      meta: {
        requiresAuth: false // 不需要登录
      }
    },

    // 主布局路由（嵌套路由）
    {
      path: '/',
      name: 'Layout',
      component: Layout,
      meta: {
        requiresAuth: false // 需要登录
      },
      children: [
        // 首页
        {
          path: '',
          name: 'Home',
          component: () => import('@/views/home/index.vue')
        },
        {
          path: 'publish',
          name: 'Publish',
          component: () => import('@/views/publish/index.vue')
        },
          // 文章列表
        {
          path: '/article',
          name: 'ArticleList',
          component: () => import('@/views/article/index.vue')
        },
        // 用户列表
        // {
        //   path: '/user',
        //   name: 'UserList',
        //   component: () => import('@/views/user/index.vue')
        // },
        // // 新增用户
        // {
        //   path: '/user/add',
        //   name: 'UserAdd',
        //   component: () => import('@/views/User/Add.vue')
        // },
        // error 页面
        {
          path: '/:pathMatch(.*)*',
          name: 'NotFound',
          component: () => import('@/views/errors/404.vue')
        }
      ]
    }
  ]
})

// 路由守卫：验证登录状态
router.beforeEach((to, from, next) => {
  // 检查路由是否需要登录
  if (to.meta.requiresAuth) {
    // 检查是否有token
    const token = localStorage.getItem('token')

    if (token) {
      // 有token，允许访问
      next()
    } else {
      // 无token，跳转到登录页，并记录当前地址以便登录后返回
      next({
        path: '/login',
        query: { redirect: to.fullPath } // 存储跳转前的路径
      })
    }
  } else {
    // 不需要登录的页面直接放行
    next()
  }
})

// 路由跳转后处理
router.afterEach((to) => {
  // 可以在这里设置页面标题等
  document.title = to.meta.title || '后台管理系统'
})

export default router
