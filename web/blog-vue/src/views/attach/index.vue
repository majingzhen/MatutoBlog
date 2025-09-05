<template>
  <div class="attach-page">
    <!-- 页面标题和操作 -->
    <div class="page-header">
      <h2>附件管理</h2>
      <el-button type="primary" @click="handleUpload">
        <el-icon><Upload /></el-icon>
        上传附件
      </el-button>
    </div>

    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="文件名">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入文件名"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 附件列表 -->
    <el-card class="table-card">
      <el-table 
        :data="attachList" 
        v-loading="loading"
        border
        stripe
      >
        <el-table-column label="ID" prop="id" width="80" />
        <el-table-column label="文件名" prop="name" min-width="200">
          <template #default="scope">
            <el-tooltip :content="scope.row.name" placement="top">
              <span class="file-name">{{ scope.row.name }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="文件路径" prop="path" min-width="250">
          <template #default="scope">
            <el-tooltip :content="scope.row.path" placement="top">
              <span class="file-path">{{ scope.row.path }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="预览" width="100">
          <template #default="scope">
            <el-button 
              type="text" 
              size="small" 
              @click="handlePreview(scope.row)"
              v-if="isImage(scope.row.name)"
            >
              <el-icon><View /></el-icon>
            </el-button>
            <span v-else class="text-muted">不支持</span>
          </template>
        </el-table-column>
        <el-table-column label="备注" prop="remark" width="150">
          <template #default="scope">
            <span>{{ scope.row.remark || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="上传时间" prop="createdAt" width="180">
          <template #default="scope">
            <span>{{ formatTime(scope.row.createdAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button
              type="text"
              size="small"
              @click="handleCopyUrl(scope.row)"
            >
              <el-icon><CopyDocument /></el-icon>
              复制链接
            </el-button>
            <el-button
              type="text"
              size="small"
              class="delete-btn"
              @click="handleDelete(scope.row)"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
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

    <!-- 上传对话框 -->
    <el-dialog
      v-model="uploadDialog"
      title="上传附件"
      width="500px"
      @close="resetUpload"
    >
      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :multiple="true"
        :file-list="selectedFiles"
        :on-change="handleFileChange"
        :on-remove="handleFileRemove"
        :show-file-list="false"
        drag
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          拖拽文件到此处或<em>点击选择文件</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持常见格式文件，单个文件不超过10MB，选择后立即上传
          </div>
        </template>
      </el-upload>

      <!-- 上传状态列表 -->
      <div v-if="uploadingFiles.length > 0" class="upload-status-list">
        <div class="status-header">
          <span>上传进度</span>
          <el-button 
            type="text" 
            size="small" 
            @click="clearCompletedUploads"
            v-if="uploadingFiles.some(f => f.status !== 'uploading')"
          >
            清除已完成
          </el-button>
        </div>
        <div class="upload-item" v-for="file in uploadingFiles" :key="file.id" :class="{ 'has-error': file.status === 'error' }">
          <div class="file-info">
            <el-icon class="file-icon">
              <Document v-if="!isImage(file.name)" />
              <Picture v-else />
            </el-icon>
            <div class="file-details">
              <div class="file-name">{{ file.name }}</div>
              <div class="file-size">{{ formatFileSize(file.size) }}</div>
            </div>
          </div>
          <div class="upload-progress">
            <div class="status-badge">
              <el-icon v-if="file.status === 'uploading'" class="uploading">
                <Loading />
              </el-icon>
              <el-icon v-else-if="file.status === 'success'" class="success">
                <Check />
              </el-icon>
              <el-icon v-else-if="file.status === 'error'" class="error">
                <Close />
              </el-icon>
            </div>
          </div>
          <div v-if="file.status === 'error'" class="error-message">
            {{ file.error }}
          </div>
        </div>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="uploadDialog = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 图片预览对话框 -->
    <el-dialog
      v-model="previewDialog"
      title="图片预览"
      width="80%"
      :append-to-body="true"
    >
      <div class="preview-container">
        <img 
          :src="previewUrl" 
          :alt="previewName"
          style="max-width: 100%; height: auto;"
        />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="previewDialog = false">关闭</el-button>
          <el-button type="primary" @click="handleDownload">下载</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, UploadFilled, View, Delete, CopyDocument, Document, Picture, Loading, Check, Close } from '@element-plus/icons-vue'
import { getAttachList, uploadAttach, deleteAttach } from '@/api/attach'

// 响应式数据
const loading = ref(false)
const attachList = ref([])
const uploadDialog = ref(false)
const previewDialog = ref(false)
const uploading = ref(false)
const uploadRef = ref()
const selectedFiles = ref([])
const previewUrl = ref('')
const previewName = ref('')
const uploadingFiles = ref([]) // 上传中的文件列表

// 搜索表单
const searchForm = reactive({
  name: ''
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  pages: 0
})

// 获取附件列表
const fetchAttachList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      name: searchForm.name || undefined
    }
    
    const response = await getAttachList(params)
    const { list, total, pages } = response.data
    
    attachList.value = list || []
    pagination.total = total || 0
    pagination.pages = pages || 0
  } catch (error) {
    console.error('获取附件列表失败:', error)
    ElMessage.error('获取附件列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchAttachList()
}

// 重置搜索
const handleReset = () => {
  searchForm.name = ''
  pagination.page = 1
  fetchAttachList()
}

// 分页大小变化
const handleSizeChange = (newSize) => {
  pagination.pageSize = newSize
  pagination.page = 1
  fetchAttachList()
}

// 当前页变化
const handleCurrentChange = (newPage) => {
  pagination.page = newPage
  fetchAttachList()
}

// 打开上传对话框
const handleUpload = () => {
  uploadDialog.value = true
}

// 文件选择变化
const handleFileChange = (file, fileList) => {
  console.log('文件选择变化:', file, fileList)
  // 文件大小限制 10MB
  if (file.raw && file.raw.size > 10 * 1024 * 1024) {
    ElMessage.error('文件大小不能超过10MB')
    return false
  }

  // 立即上传文件
  if (file.raw) {
    handleImmediateUpload(file)
  }

  return true
}

// 移除文件
const handleFileRemove = (file, fileList) => {
  // Element Plus 会自动更新 fileList
}

// 立即上传单个文件
const handleImmediateUpload = async (fileItem) => {
  const fileObj = {
    id: Date.now() + Math.random(),
    name: fileItem.name,
    size: fileItem.size,
    status: 'uploading', // uploading, success, error
    response: null,
    error: null
  }
  
  uploadingFiles.value.push(fileObj)
  
  try {
    const response = await uploadAttach(fileItem.raw)
    
    // 更新文件状态
    const index = uploadingFiles.value.findIndex(f => f.id === fileObj.id)
    if (index !== -1) {
      uploadingFiles.value[index].status = 'success'
      uploadingFiles.value[index].response = response.data
    }
    
    ElMessage.success(`文件 ${fileItem.name} 上传成功`)
    fetchAttachList() // 重新加载列表
    
    // 3秒后自动清除成功的上传记录
    setTimeout(() => {
      const successIndex = uploadingFiles.value.findIndex(f => f.id === fileObj.id && f.status === 'success')
      if (successIndex !== -1) {
        uploadingFiles.value.splice(successIndex, 1)
      }
    }, 3000)
  } catch (error) {
    console.error('上传失败:', error)
    
    // 更新文件状态
    const index = uploadingFiles.value.findIndex(f => f.id === fileObj.id)
    if (index !== -1) {
      uploadingFiles.value[index].status = 'error'
      uploadingFiles.value[index].error = error.message || '上传失败'
    }
    
    ElMessage.error(`文件 ${fileItem.name} 上传失败`)
  }
}

// 确认上传（保留用于批量上传的备用功能）
const handleConfirmUpload = async () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  uploading.value = true
  
  try {
    // 逐个上传文件
    const uploadPromises = selectedFiles.value.map(fileItem => {
      return handleImmediateUpload(fileItem)
    })
    
    await Promise.allSettled(uploadPromises)
    
    uploadDialog.value = false
    resetUpload()
  } catch (error) {
    console.error('上传过程出错:', error)
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

// 重置上传状态
const resetUpload = () => {
  selectedFiles.value = []
  uploading.value = false
  uploadingFiles.value = [] // 清空上传状态列表
  if (uploadRef.value) {
    uploadRef.value.clearFiles()
  }
}

// 清除已完成的上传记录
const clearCompletedUploads = () => {
  uploadingFiles.value = uploadingFiles.value.filter(file => file.status === 'uploading')
}

// 判断是否为图片文件
const isImage = (fileName) => {
  const imageExt = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp']
  const ext = fileName.toLowerCase().split('.').pop()
  return imageExt.includes(ext)
}

// 预览图片
const handlePreview = (row) => {
  previewUrl.value = `http://localhost:8080/uploads/${row.path}`
  previewName.value = row.name
  previewDialog.value = true
}

// 下载文件
const handleDownload = () => {
  const link = document.createElement('a')
  link.href = previewUrl.value
  link.download = previewName.value
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 复制文件链接
const handleCopyUrl = async (row) => {
  const url = `http://localhost:8080/uploads/${row.path}`
  
  try {
    await navigator.clipboard.writeText(url)
    ElMessage.success('链接已复制到剪贴板')
  } catch (error) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = url
    textArea.style.position = 'fixed'
    textArea.style.opacity = '0'
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    ElMessage.success('链接已复制到剪贴板')
  }
}

// 删除附件
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除附件「${row.name}」吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteAttach(row.id)
      ElMessage.success('删除成功')
      fetchAttachList() // 重新加载列表
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    // 用户取消删除
  })
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return '-'
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 格式化文件大小
const formatFileSize = (size) => {
  if (!size) return '-'
  const units = ['B', 'KB', 'MB', 'GB']
  let index = 0
  while (size >= 1024 && index < units.length - 1) {
    size /= 1024
    index++
  }
  return `${size.toFixed(1)} ${units[index]}`
}

// 页面加载时获取数据
onMounted(() => {
  fetchAttachList()
})
</script>

<style scoped>
.attach-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  color: #303133;
  font-size: 24px;
  font-weight: 500;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.file-name, 
.file-path {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.text-muted {
  color: #909399;
  font-size: 12px;
}

.delete-btn {
  color: #f56c6c;
}

.delete-btn:hover {
  color: #f56c6c;
  background-color: #fef0f0;
}

.pagination-container {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.preview-container {
  text-align: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.el-upload--drag {
  width: 100%;
  transition: all 0.3s ease;
}

.el-upload--drag:hover {
  border-color: #409eff;
  background-color: #f0f8ff;
}

.el-upload--drag.is-dragover {
  border-color: #409eff;
  background-color: #e6f3ff;
  border-style: solid;
  animation: dragover-pulse 1s infinite alternate;
}

@keyframes dragover-pulse {
  0% {
    transform: scale(1);
    box-shadow: 0 0 0 0 rgba(64, 158, 255, 0.4);
  }
  100% {
    transform: scale(1.02);
    box-shadow: 0 0 0 10px rgba(64, 158, 255, 0.1);
  }
}

.el-upload__tip {
  margin-top: 10px;
  color: #606266;
  font-size: 13px;
  line-height: 1.5;
}

/* 上传状态列表样式 */
.upload-status-list {
  margin-top: 20px;
  max-height: 300px;
  overflow-y: auto;
}

.status-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-weight: 500;
  color: #303133;
}

.upload-item {
  display: flex;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  background-color: #f8f9fa;
  border-radius: 6px;
  border: 1px solid #ebeef5;
  position: relative;
}

.upload-item.has-error {
  margin-bottom: 28px;
}

.file-info {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
}

.file-icon {
  font-size: 20px;
  margin-right: 8px;
  color: #909399;
}

.file-details {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 14px;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-size {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.upload-progress {
  display: flex;
  align-items: center;
  margin-left: 12px;
  min-width: 40px;
}

.status-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
}

.status-badge .uploading {
  color: #409eff;
  animation: rotating 2s linear infinite;
}

.status-badge .success {
  color: #67c23a;
}

.status-badge .error {
  color: #f56c6c;
}

.error-message {
  position: absolute;
  bottom: -20px;
  left: 0;
  font-size: 12px;
  color: #f56c6c;
}

@keyframes rotating {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .attach-page {
    padding: 10px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .el-table {
    font-size: 12px;
  }
}
</style>