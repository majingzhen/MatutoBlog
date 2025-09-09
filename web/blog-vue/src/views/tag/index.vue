<template>
  <div class="tag-list-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">标签管理</h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增标签
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
        <el-form-item label="标签名称">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入标签名称"
            clearable
            style="width: 200px"
            @clear="handleSearch"
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

    <!-- 标签列表表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="tagList"
        style="width: 100%"
        stripe
        border
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column prop="name" label="标签名称" min-width="150">
          <template #default="{ row }">
            <div class="tag-name">
              <el-tag 
                :color="row.color || '#007bff'"
                effect="light"
                style="margin-right: 8px"
              >
                {{ row.name }}
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="color" label="标签颜色" width="120">
          <template #default="{ row }">
            <div class="color-preview" :style="{ backgroundColor: row.color || '#007bff' }"></div>
            <span style="margin-left: 8px">{{ row.color || '#007bff' }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="article_count" label="文章数" width="100">
          <template #default="{ row }">
            <el-tag type="info" size="small">
              {{ row.article_count || 0 }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="sort" label="排序" width="80">
          <template #default="{ row }">
            <span>{{ row.sort || 0 }}</span>
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

    <!-- 新增/编辑标签对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑标签' : '新增标签'"
      width="600px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="80px"
      >
        <el-form-item label="标签名称" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入标签名称"
            maxlength="20"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="标签颜色" prop="color">
          <div class="color-picker-wrapper">
            <el-color-picker
              v-model="form.color"
              show-alpha
              :predefine="predefineColors"
            />
            <el-input
              v-model="form.color"
              placeholder="请选择或输入颜色值"
              style="width: 200px; margin-left: 15px"
            />
          </div>
        </el-form-item>
        
        <el-form-item label="排序" prop="sort">
          <el-input-number
            v-model="form.sort"
            :min="0"
            :max="999"
            controls-position="right"
            placeholder="数字越小排序越靠前"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            :loading="submitLoading" 
            @click="handleSubmit"
          >
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, Search, RefreshRight, Edit, Delete 
} from '@element-plus/icons-vue'
import { getTagList, createTag, updateTag, deleteTag } from '@/api/tag.js'

// 响应式数据
const loading = ref(false)
const submitLoading = ref(false)
const tagList = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

// 搜索表单
const searchForm = reactive({
  name: ''
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 表单数据
const form = reactive({
  id: null,
  name: '',
  color: '#007bff',
  sort: 0
})

// 预定义颜色
const predefineColors = [
  '#ff4500',
  '#ff8c00',
  '#ffd700',
  '#90ee90',
  '#00ced1',
  '#1e90ff',
  '#c71585',
  '#007bff',
  '#28a745',
  '#dc3545',
  '#6c757d',
  '#17a2b8'
]

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' },
    { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' }
  ],
  color: [
    { required: true, message: '请选择标签颜色', trigger: 'blur' }
  ],
  sort: [
    { required: true, message: '请输入排序值', trigger: 'blur' },
    { type: 'number', min: 0, max: 999, message: '排序值应在 0 到 999 之间', trigger: 'blur' }
  ]
}

// 获取标签列表
const fetchTagList = async () => {
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

    const response = await getTagList(params)
    console.log('获取标签列表:', response)
    if (response.code == 200) {
      tagList.value = response.data.list || []
      pagination.total = response.data.total || 0
    } else {
      ElMessage.error(response.data.message || '获取标签列表失败')
    }
  } catch (error) {
    console.error('获取标签列表失败:', error)
    ElMessage.error('获取标签列表失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  fetchTagList()
}

// 重置搜索
const handleReset = () => {
  searchForm.name = ''
  pagination.page = 1
  fetchTagList()
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchTagList()
}

// 页码变化
const handleCurrentChange = (page) => {
  pagination.page = page
  fetchTagList()
}

// 新增标签
const handleCreate = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
}

// 编辑标签
const handleEdit = (row) => {
  isEdit.value = true
  dialogVisible.value = true
  
  // 填充表单数据
  form.id = row.id
  form.name = row.name
  form.color = row.color || '#007bff'
  form.sort = row.sort || 0
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    const valid = await formRef.value.validate()
    if (!valid) return
    
    submitLoading.value = true
    
    const formData = {
      name: form.name,
      color: form.color,
      sort: form.sort
    }
    
    let response
    if (isEdit.value) {
      response = await updateTag(form.id, formData)
    } else {
      response = await createTag(formData)
    }
    
    console.log('提交标签:', response)
    if (response.code == 200) {
      ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
      dialogVisible.value = false
      fetchTagList()
    } else {
      ElMessage.error(response.data.message || '操作失败')
    }
  } catch (error) {
    console.error('提交标签失败:', error)
    ElMessage.error('操作失败，请重试')
  } finally {
    submitLoading.value = false
  }
}

// 删除标签
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除标签 "${row.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    loading.value = true
    const response = await deleteTag(row.id)
    console.log('删除标签:', response)
    if (response.code == 200) {
      ElMessage.success('删除成功')
      fetchTagList()
    } else {
      ElMessage.error(response.data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除标签失败:', error)
      ElMessage.error('删除失败，请重试')
    }
  } finally {
    loading.value = false
  }
}

// 重置表单
const resetForm = () => {
  form.id = null
  form.name = ''
  form.color = '#007bff'
  form.sort = 0
  
  if (formRef.value) {
    formRef.value.clearValidate()
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
  fetchTagList()
})
</script>

<style scoped>
.tag-list-page {
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

.tag-name {
  display: flex;
  align-items: center;
}

.color-preview {
  display: inline-block;
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 1px solid #dcdfe6;
  vertical-align: middle;
}

.color-picker-wrapper {
  display: flex;
  align-items: center;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .tag-list-page {
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
  
  .color-picker-wrapper {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
}
</style>