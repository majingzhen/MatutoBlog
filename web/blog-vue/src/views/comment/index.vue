<template>
  <div class="comment-list-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">评论管理</h1>
      <div class="header-actions">
        <el-button 
          type="success" 
          :disabled="selectedComments.length === 0"
          @click="handleBatchReview(1)"
        >
          <el-icon><Check /></el-icon>
          批量通过
        </el-button>
        <el-button 
          type="warning" 
          :disabled="selectedComments.length === 0"
          @click="handleBatchReview(2)"
        >
          <el-icon><Close /></el-icon>
          批量拒绝
        </el-button>
      </div>
    </div>

    <!-- 搜索筛选区域 -->
    <el-card class="search-card" shadow="never">
      <el-form 
        :model="searchForm" 
        :inline="true" 
        class="search-form"
        @submit.prevent="handleSearch"
      >
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索用户名、内容或邮箱"
            clearable
            style="width: 250px"
            @clear="handleSearch"
          />
        </el-form-item>

        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
            @change="handleSearch"
          >
            <el-option label="正常" :value="0" />
            <el-option label="待审核" :value="1" />
            <el-option label="已拒绝" :value="2" />
          </el-select>
        </el-form-item>

        <el-form-item label="文章ID">
          <el-input-number
            v-model="searchForm.articleId"
            placeholder="文章ID"
            :min="1"
            controls-position="right"
            style="width: 120px"
            @change="handleSearch"
          />
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

    <!-- 评论列表表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="commentList"
        style="width: 100%"
        stripe
        border
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column prop="username" label="评论人" width="120">
          <template #default="{ row }">
            <div class="user-info">
              <div>{{ row.username }}</div>
              <small class="text-muted">{{ row.email }}</small>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="content" label="评论内容" min-width="300">
          <template #default="{ row }">
            <div class="comment-content">
              <p>{{ row.content }}</p>
              <div class="comment-meta">
                <el-tag v-if="row.article" size="small" type="info">
                  文章: {{ row.article.title }}
                </el-tag>
                <span class="ip-info">IP: {{ row.ip }}</span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag 
              :type="getStatusType(row.status)" 
              size="small"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="website" label="网站" width="120">
          <template #default="{ row }">
            <el-link 
              v-if="row.website" 
              :href="row.website" 
              target="_blank" 
              type="primary"
              :underline="false"
            >
              访问网站
            </el-link>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="评论时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button 
                v-if="row.status !== 1"
                type="success" 
                size="small" 
                @click="handleReview(row, 1)"
              >
                <el-icon><Check /></el-icon>
                通过
              </el-button>
              <el-button 
                v-if="row.status !== 2"
                type="warning" 
                size="small" 
                @click="handleReview(row, 2)"
              >
                <el-icon><Close /></el-icon>
                拒绝
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Check, Close, Search, RefreshRight, Delete 
} from '@element-plus/icons-vue'
import { 
  getCommentList, 
  updateCommentStatus, 
  deleteComment, 
  batchReviewComments 
} from '@/api/comment.js'

// 响应式数据
const loading = ref(false)
const commentList = ref([])
const selectedComments = ref([])

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: null,
  articleId: null
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 状态类型映射
const getStatusType = (status) => {
  const statusMap = {
    0: 'success',  // 正常
    1: 'warning',  // 待审核
    2: 'danger'    // 已拒绝
  }
  return statusMap[status] || 'info'
}

// 状态文本映射
const getStatusText = (status) => {
  const statusMap = {
    0: '正常',
    1: '待审核',
    2: '已拒绝'
  }
  return statusMap[status] || '未知'
}

// 获取评论列表
const fetchCommentList = async () => {
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

    const response = await getCommentList(params)
    console.log('获取评论列表:', response)
    if (response.code == 200) {
      commentList.value = response.data.list || []
      pagination.total = response.data.total || 0
    } else {
      ElMessage.error(response.data.message || '获取评论列表失败')
    }
  } catch (error) {
    console.error('获取评论列表失败:', error)
    ElMessage.error('获取评论列表失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  fetchCommentList()
}

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = null
  searchForm.articleId = null
  pagination.page = 1
  fetchCommentList()
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchCommentList()
}

// 页码变化
const handleCurrentChange = (page) => {
  pagination.page = page
  fetchCommentList()
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedComments.value = selection
}

// 审核评论
const handleReview = async (row, status) => {
  try {
    const statusText = getStatusText(status)
    await ElMessageBox.confirm(
      `确定要将评论状态设置为"${statusText}"吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    loading.value = true
    const response = await updateCommentStatus(row.id, status)
    console.log('审核评论:', response)
    if (response.code == 200) {
      ElMessage.success('操作成功')
      fetchCommentList()
    } else {
      ElMessage.error(response.data.message || '操作失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('审核评论失败:', error)
      ElMessage.error('操作失败，请重试')
    }
  } finally {
    loading.value = false
  }
}

// 批量审核
const handleBatchReview = async (status) => {
  if (selectedComments.value.length === 0) {
    ElMessage.warning('请先选择要操作的评论')
    return
  }

  try {
    const statusText = getStatusText(status)
    await ElMessageBox.confirm(
      `确定要将选中的 ${selectedComments.value.length} 条评论状态设置为"${statusText}"吗？`,
      '确认批量操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const ids = selectedComments.value.map(item => item.id)
    loading.value = true
    const response = await batchReviewComments(ids, status)
    console.log('批量审核评论:', response)
    if (response.code == 200) {
      ElMessage.success('批量操作成功')
      fetchCommentList()
      selectedComments.value = []
    } else {
      ElMessage.error(response.data.message || '批量操作失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量审核失败:', error)
      ElMessage.error('批量操作失败，请重试')
    }
  } finally {
    loading.value = false
  }
}

// 删除评论
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户"${row.username}"的这条评论吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    loading.value = true
    const response = await deleteComment(row.id)
    console.log('删除评论:', response)
    if (response.code == 200) {
      ElMessage.success('删除成功')
      fetchCommentList()
    } else {
      ElMessage.error(response.data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除评论失败:', error)
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
  fetchCommentList()
})
</script>

<style scoped>
.comment-list-page {
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

.header-actions {
  display: flex;
  gap: 10px;
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

.user-info {
  line-height: 1.4;
}

.text-muted {
  color: #909399;
  font-size: 12px;
}

.comment-content p {
  margin: 0 0 8px 0;
  line-height: 1.5;
  max-height: 4.5em;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.comment-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 8px;
}

.ip-info {
  font-size: 12px;
  color: #909399;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .comment-list-page {
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

  .comment-content p {
    -webkit-line-clamp: 2;
    max-height: 3em;
  }
}
</style>