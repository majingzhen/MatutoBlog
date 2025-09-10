<template>
  <div class="article-list-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">文章管理</h1>
      <el-button type="primary" @click="handleCreate">
        <Icon icon="material-symbols:add" />
        新增文章
      </el-button>
    </div>

    <!-- 搜索筛选区域 -->
    <el-card class="search-card" shadow="never">
      <el-form 
        :model="searchForm" 
        :inline="true" 
        class="search-form"
        @submit.prevent="handleSearch"
      >
        <el-form-item label="文章标题">
          <el-input
            v-model="searchForm.title"
            placeholder="请输入文章标题"
            clearable
            style="width: 200px"
            @clear="handleSearch"
          />
        </el-form-item>
        
        <el-form-item label="分类">
          <el-select
            v-model="searchForm.categoryId"
            placeholder="请选择分类"
            clearable
            style="width: 150px"
            @change="handleSearch"
          >
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
            @change="handleSearch"
          >
            <el-option label="已发布" :value="0" />
            <el-option label="草稿" :value="1" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><RefreshRight /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 文章列表表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="articleList"
        style="width: 100%"
        stripe
        border
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column prop="title" label="文章标题" min-width="200">
          <template #default="{ row }">
            <div class="article-title">
              <el-link 
                type="primary" 
                :underline="false"
                @click="handleView(row)"
              >
                {{ row.title }}
              </el-link>
              <el-tag v-if="row.is_top" type="warning" size="small" style="margin-left: 8px">
                置顶
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.category" type="info" size="small">
              {{ row.category.name }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'info'" size="small">
              {{ row.status === 0 ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="view_count" label="阅读量" width="100">
          <template #default="{ row }">
            <span>{{ row.view_count || 0 }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button 
                type="primary" 
                size="small" 
                @click="handleEdit(row)"
              >
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button 
                type="danger" 
                size="small" 
                @click="handleDelete(row)"
              >
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页组件 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, Search, RefreshRight, Edit, Delete 
} from '@element-plus/icons-vue'
import axios from 'axios'
import { getArticleList, deleteArticle } from '@/api/article.js'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const articleList = ref([])

// 搜索表单
const searchForm = reactive({
  title: '',
  categoryId: null,
  status: null
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 分类数据
const categories = ref([
  { id: 1, name: '前端开发' },
  { id: 2, name: '后端开发' },
  { id: 3, name: '数据库' },
  { id: 4, name: '服务器' },
  { id: 5, name: '人工智能' }
])

// 获取文章列表
const fetchArticleList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchForm
    }
    
    // 清除空值参数
    Object.keys(params).forEach(key => {
      if (params[key] === null || params[key] === '') {
        delete params[key]
      }
    })

    const response = await getArticleList(params)
    console.log('获取文章列表:', response)
    if (response.code == 200) {
      articleList.value = response.data.list || []
      pagination.total = response.data.total || 0
    } else {
      ElMessage.error(response.data.message || '获取文章列表失败')
    }
  } catch (error) {
    console.error('获取文章列表失败:', error)
    ElMessage.error('获取文章列表失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  fetchArticleList()
}

// 重置搜索
const handleReset = () => {
  searchForm.title = ''
  searchForm.categoryId = null
  searchForm.status = null
  pagination.page = 1
  fetchArticleList()
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchArticleList()
}

// 页码变化
const handleCurrentChange = (page) => {
  pagination.page = page
  fetchArticleList()
}

// 新增文章
const handleCreate = () => {
  router.push('/publish')
}

// 查看文章
const handleView = (row) => {
  window.open(`http://127.0.0.1:8080/article/${row.id}`, '_blank')
}

// 编辑文章
const handleEdit = (row) => {
  router.push(`/publish?id=${row.id}`)
}

// 删除文章
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文章 "${row.title}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    loading.value = true
    const response = await deleteArticle(row.id)
    console.log('删除文章:', response)
    if (response.code == 200) {
      ElMessage.success('删除成功')
      fetchArticleList()
    } else {
      ElMessage.error(response.data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除文章失败:', error)
      ElMessage.error('删除失败，请重试')
    }
  } finally {
    loading.value = false
  }
}

// 格式化日期时间
const formatDateTime = (dateTime) => {
  if (!dateTime) return '-'
  const date = new Date(dateTime)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 组件挂载时获取数据
onMounted(() => {
  fetchArticleList()
})
</script>

<style scoped>
.article-list-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #1d2129;
}

.search-card {
  margin-bottom: 20px;
}

.search-form {
  margin-bottom: -10px;
}

.table-card {
  min-height: 400px;
}

.article-title {
  display: flex;
  align-items: center;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .article-list-page {
    padding: 10px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .search-form {
    flex-direction: column;
  }
  
  .search-form .el-form-item {
    margin-right: 0;
    margin-bottom: 15px;
  }
}
</style>