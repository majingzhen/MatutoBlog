<template>
  <div class="category-list-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">分类管理</h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新增分类
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
        <el-form-item label="分类名称">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入分类名称"
            clearable
            style="width: 200px"
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
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
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

    <!-- 分类列表表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="categoryList"
        style="width: 100%"
        stripe
        border
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column prop="name" label="分类名称" min-width="150">
          <template #default="{ row }">
            <div class="category-name">
              <el-link 
                type="primary" 
                :underline="false"
              >
                {{ row.name }}
              </el-link>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="desc" label="描述" min-width="200">
          <template #default="{ row }">
            <span>{{ row.desc || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="slug" label="Slug" width="120">
          <template #default="{ row }">
            <el-tag type="info" size="small">
              {{ row.slug || '-' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="pid" label="父分类ID" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.pid !== -1" type="warning" size="small">
              {{ row.pid }}
            </el-tag>
            <el-tag v-else type="success" size="small">
              根分类
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
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

    <!-- 新增/编辑分类对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑分类' : '新增分类'"
      width="600px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="80px"
      >
        <el-form-item label="分类名称" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入分类名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="父分类" prop="pid">
          <el-input-number
            v-model="form.pid"
            :min="-1"
            controls-position="right"
            placeholder="-1为根分类"
          />
        </el-form-item>
        
        <el-form-item label="分类描述" prop="desc">
          <el-input
            v-model="form.desc"
            type="textarea"
            :rows="3"
            placeholder="请输入分类描述"
            maxlength="512"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="缩略图" prop="thumbnail">
          <el-input
            v-model="form.thumbnail"
            placeholder="请输入缩略图URL"
            maxlength="256"
          />
        </el-form-item>
        
        <el-form-item label="Slug" prop="slug">
          <el-input
            v-model="form.slug"
            placeholder="URL别名，留空自动生成"
            maxlength="128"
          />
        </el-form-item>
        
        <el-form-item label="SEO关键词" prop="metaKeywords">
          <el-input
            v-model="form.metaKeywords"
            placeholder="SEO关键词，多个用逗号分隔"
            maxlength="256"
          />
        </el-form-item>
        
        <el-form-item label="SEO描述" prop="metaDescription">
          <el-input
            v-model="form.metaDescription"
            type="textarea"
            :rows="2"
            placeholder="SEO描述内容"
            maxlength="256"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
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
import { getCategoryList, createCategory, updateCategory, deleteCategory } from '@/api/category.js'

// 响应式数据
const loading = ref(false)
const submitLoading = ref(false)
const categoryList = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

// 搜索表单
const searchForm = reactive({
  name: '',
  status: null
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
  pid: -1,
  desc: '',
  thumbnail: '',
  slug: '',
  metaKeywords: '',
  metaDescription: '',
  status: 0
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 1, max: 256, message: '长度在 1 到 256 个字符', trigger: 'blur' }
  ],
  pid: [
    { required: true, message: '请输入父分类ID', trigger: 'blur' },
    { type: 'number', min: -1, message: '父分类ID不能小于-1', trigger: 'blur' }
  ],
  desc: [
    { max: 512, message: '描述不能超过 512 个字符', trigger: 'blur' }
  ],
  thumbnail: [
    { max: 256, message: 'URL长度不能超过 256 个字符', trigger: 'blur' }
  ],
  slug: [
    { max: 128, message: 'Slug长度不能超过 128 个字符', trigger: 'blur' }
  ],
  metaKeywords: [
    { max: 256, message: 'SEO关键词长度不能超过 256 个字符', trigger: 'blur' }
  ],
  metaDescription: [
    { max: 256, message: 'SEO描述长度不能超过 256 个字符', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 获取分类列表
const fetchCategoryList = async () => {
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

    const response = await getCategoryList(params)
    console.log('获取分类列表:', response)
    if (response.code == 200) {
      categoryList.value = response.data.list || []
      pagination.total = response.data.total || 0
    } else {
      ElMessage.error(response.data.message || '获取分类列表失败')
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
    ElMessage.error('获取分类列表失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  fetchCategoryList()
}

// 重置搜索
const handleReset = () => {
  searchForm.name = ''
  searchForm.status = null
  pagination.page = 1
  fetchCategoryList()
}

// 分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchCategoryList()
}

// 页码变化
const handleCurrentChange = (page) => {
  pagination.page = page
  fetchCategoryList()
}

// 新增分类
const handleCreate = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
}

// 编辑分类
const handleEdit = (row) => {
  isEdit.value = true
  dialogVisible.value = true
  
  // 填充表单数据
  form.id = row.id
  form.name = row.name
  form.pid = row.pid || -1
  form.desc = row.desc || ''
  form.thumbnail = row.thumbnail || ''
  form.slug = row.slug || ''
  form.metaKeywords = row.metaKeywords || ''
  form.metaDescription = row.metaDescription || ''
  form.status = row.status
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
      pid: form.pid,
      desc: form.desc,
      thumbnail: form.thumbnail,
      slug: form.slug,
      metaKeywords: form.metaKeywords,
      metaDescription: form.metaDescription,
      status: form.status
    }
    
    let response
    if (isEdit.value) {
      response = await updateCategory(form.id, formData)
    } else {
      response = await createCategory(formData)
    }
    
    console.log('提交分类:', response)
    if (response.code == 200) {
      ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
      dialogVisible.value = false
      fetchCategoryList()
    } else {
      ElMessage.error(response.data.message || '操作失败')
    }
  } catch (error) {
    console.error('提交分类失败:', error)
    ElMessage.error('操作失败，请重试')
  } finally {
    submitLoading.value = false
  }
}

// 删除分类
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除分类 "${row.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    loading.value = true
    const response = await deleteCategory(row.id)
    console.log('删除分类:', response)
    if (response.code == 200) {
      ElMessage.success('删除成功')
      fetchCategoryList()
    } else {
      ElMessage.error(response.data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除分类失败:', error)
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
  form.pid = -1
  form.desc = ''
  form.thumbnail = ''
  form.slug = ''
  form.metaKeywords = ''
  form.metaDescription = ''
  form.status = 0
  
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
  fetchCategoryList()
})
</script>

<style scoped>
.category-list-page {
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

.category-name {
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
  .category-list-page {
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