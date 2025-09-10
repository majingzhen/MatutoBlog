<template>
  <div class="article-editor-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <el-page-header @back="handleBack" :content="isEdit ? '编辑文章' : '新建文章'" />
      <div class="header-actions">
        <el-button @click="handlePreview" type="info" plain>
          <Icon icon="mdi:eye" />
          预览
        </el-button>
        <el-button 
          size="default"
          @click="handleSaveDraft"
          :loading="saving"
        >
          <el-icon><Save /></el-icon>
          保存草稿
        </el-button>
        <el-button 
          type="primary" 
          size="default"
          @click="handlePublish"
          :loading="submitting"
        >
          <el-icon><Upload /></el-icon>
          {{ isEdit ? '更新文章' : '发布文章' }}
        </el-button>
      </div>
    </div>

    <div class="editor-container">
      <!-- 左侧编辑区域 -->
      <div class="editor-main">
        <!-- 标题输入框 -->
        <div class="title-input-wrapper">
          <el-input
            v-model="article.title"
            placeholder="请输入文章标题"
            size="large"
            :maxlength="200"
            show-word-limit
            class="title-input"
            @blur="generateSlug"
          />
        </div>

        <!-- Vditor编辑器 -->
        <div class="editor-wrapper">
          <div id="vditor" class="vditor-container"></div>
        </div>
      </div>

      <!-- 右侧配置区域 -->
      <div class="editor-sidebar">
        <!-- 文章配置表单 -->
        <div class="article-config-form">
          <div class="config-form-content">
            <!-- 访问地址别名 -->
            <div class="form-item">
              <label class="form-label">访问地址别名</label>
              <el-input
                v-model="article.slug"
                placeholder="请输入访问地址别名"
                :maxlength="100"
                size="default"
              />
            </div>

            <!-- 分类 -->
            <div class="form-item">
              <label class="form-label">分类</label>
              <el-select
                v-model="article.categoryId"
                placeholder="请选择分类"
                clearable
                style="width: 100%;"
              >
                <el-option
                  v-for="category in categories"
                  :key="category.id"
                  :label="category.name"
                  :value="category.id"
                />
              </el-select>
            </div>

            <!-- 标签 -->
            <div class="form-item">
              <label class="form-label">标签</label>
              <el-select
                v-model="selectedTags"
                multiple
                placeholder="请选择或输入新标签"
                style="width: 100%;"
                filterable
                allow-create
                :reserve-keyword="false"
                default-first-option
                @change="handleTagChange"
              >
                <el-option
                  v-for="tag in allTags"
                  :key="tag.id"
                  :label="tag.name"
                  :value="tag.name"
                />
              </el-select>
              <div v-if="selectedTags.length > 0" class="selected-tags-display">
                <el-tag 
                  v-for="tagName in selectedTags" 
                  :key="tagName" 
                  class="tag-item-display"
                  :type="isExistingTag(tagName) ? 'success' : 'info'"
                  size="small"
                  closable 
                  @close="removeTag(tagName)"
                >
                  {{ tagName }}
                  <span v-if="!isExistingTag(tagName)" class="new-tag-indicator"> (新)</span>
                </el-tag>
              </div>
            </div>

            <!-- 文章摘要 -->
            <div class="form-item">
              <label class="form-label">文章摘要</label>
              <el-input
                v-model="article.summary"
                type="textarea"
                :rows="4"
                placeholder="请输入文章摘要"
                :maxlength="500"
                show-word-limit
              />
            </div>

            <!-- SEO关键字 -->
            <div class="form-item">
              <label class="form-label">SEO关键字</label>
              <el-input
                v-model="article.metaKeywords"
                placeholder="请输入SEO关键字"
                :maxlength="200"
                show-word-limit
              />
            </div>

            <!-- SEO描述 -->
            <div class="form-item">
              <label class="form-label">SEO描述</label>
              <el-input
                v-model="article.metaDescription"
                type="textarea"
                :rows="3"
                placeholder="请输入SEO描述"
                :maxlength="160"
                show-word-limit
              />
            </div>

            <!-- 文章标识 -->
            <div class="form-item">
              <label class="form-label">文章标识</label>
              <el-input
                v-model="article.metaTitle"
                placeholder="请输入文章标识"
                :maxlength="100"
              />
            </div>

            <!-- 开关选项 -->
            <div class="form-item switch-section">
              <div class="switch-item">
                <span>
                  <span class="switch-label">允许评价</span>
                  <el-switch v-model="article.isComment" />
                </span>

                <span>
                  <span class="switch-label">置顶</span>
                  <el-switch v-model="article.isTop" />
                </span>
               <span>
                  <span class="switch-label">可见</span>
                  <el-switch v-model="article.isVisible" />
               </span>
              </div>
            </div>

            <!-- 封面图 -->
            <div class="form-item upload-section">
              <label class="form-label">封面图</label>
              <div class="upload-area">
                <el-upload
                  v-model:file-list="fileList"
                  :action="uploadUrl"
                  :headers="uploadHeaders"
                  :on-success="handleThumbnailSuccess"
                  :on-error="handleUploadError"
                  :before-upload="beforeThumbnailUpload"
                  list-type="picture-card"
                  :limit="1"
                  accept="image/*"
                  class="thumbnail-upload"
                >
                  <el-icon><Plus /></el-icon>
                  <template #file="{ file }">
                    <div>
                      <img class="el-upload-list__item-thumbnail" :src="file.url" alt="" />
                      <span class="el-upload-list__item-actions">
                        <span class="el-upload-list__item-preview" @click="handlePictureCardPreview(file)">
                          <el-icon><zoom-in /></el-icon>
                        </span>
                        <span class="el-upload-list__item-delete" @click="handleRemove(file)">
                          <el-icon><Delete /></el-icon>
                        </span>
                      </span>
                    </div>
                  </template>
                </el-upload>
                <span class="upload-tip">请选择上传封面图片</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片预览对话框 -->
    <el-dialog v-model="previewVisible" title="图片预览">
      <img w-full :src="previewUrl" alt="预览" style="width: 100%" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  View, Document, Folder, Picture, Search, Setting,
  Plus, Delete, ZoomIn, Upload
} from '@element-plus/icons-vue'

// 导入 Vditor
import Vditor from 'vditor'
import 'vditor/dist/index.css'

// 导入API
import { 
  createArticle, updateArticle, getArticleById,
  getCategoryList, getTagList
} from '@/api/article.js'

const router = useRouter()
const route = useRoute()

// 判断是否为编辑模式
const isEdit = computed(() => !!route.query.id)

// 状态控制
const submitting = ref(false)
const saving = ref(false)
const loading = ref(false)

// Vditor 编辑器实例
const vditor = ref(null)

// 文章数据
const article = reactive({
  id: null,
  title: '',
  slug: '',
  summary: '',
  content: '',
  parseContent: '',
  contentModel: 'markdown',
  type: 'article',
  categoryId: null,
  tagIds: [],
  addTags: [], // 新增标签名称数组
  thumbnail: '',
  metaTitle: '',
  metaDescription: '',
  metaKeywords: '',
  isTop: false,
  isComment: true,
  isVisible: true, // 新增可见性字段
  status: 0 // 0: 草稿, 1: 发布
})

// 分类和标签数据
const categories = ref([])
const allTags = ref([])
const selectedTags = ref([]) // 选中的标签名称数组

// 上传相关
const uploadUrl = import.meta.env.VITE_APP_UPLOAD_URL || '/api/upload'
const uploadHeaders = {
  'Authorization': localStorage.getItem('token') ? `Bearer ${localStorage.getItem('token')}` : ''
}
const fileList = ref([])
const previewVisible = ref(false)
const previewUrl = ref('')

// 渲染Markdown内容
const renderMarkdown = () => {
  // Vditor 会自动处理内容渲染，这里不需要额外处理
}

// 处理内容变化
const handleContentChange = () => {
  // Vditor 会自动处理内容变化
}

// 初始化 Vditor
const initVditor = () => {
  vditor.value = new Vditor('vditor', {
    height: 750, // 进一步增加编辑器高度
    mode: 'ir', // 即时渲染模式，无右侧预览
    placeholder: '请输入文章内容...',
    theme: 'classic',
    icon: 'material',
    upload: {
      url: uploadUrl,
      headers: uploadHeaders,
      accept: 'image/*',
      success: (editor, response) => {
        if (response.code === 200) {
          ElMessage.success('图片上传成功')
        } else {
          ElMessage.error(response.message || '上传失败')
        }
      },
      error: (message) => {
        ElMessage.error('图片上传失败: ' + message)
      }
    },
    after: () => {
      // 编辑器初始化完成后设置内容
      if (article.content) {
        vditor.value.setValue(article.content)
      }
    },
    input: (value) => {
      // 内容变化时更新数据
      article.content = value
      article.parseContent = vditor.value.getHTML()
    }
  })
}

// 标签处理相关方法
const isExistingTag = (tagName) => {
  return allTags.value.some(tag => tag.name === tagName)
}

const handleTagChange = (selectedTagNames) => {
  selectedTags.value = selectedTagNames
  // 分离已有标签ID和新标签名称
  const existingTagIds = []
  const newTagNames = []
  
  selectedTagNames.forEach(tagName => {
    const existingTag = allTags.value.find(tag => tag.name === tagName)
    if (existingTag) {
      existingTagIds.push(existingTag.id)
    } else {
      newTagNames.push(tagName)
    }
  })
  
  article.tagIds = existingTagIds
  article.addTags = newTagNames
}

const removeTag = (tagName) => {
  selectedTags.value = selectedTags.value.filter(name => name !== tagName)
  handleTagChange(selectedTags.value)
}

// 初始化数据
const initData = async () => {
  loading.value = true
  try {
    // 获取分类和标签列表
    const [categoryRes, tagRes] = await Promise.all([
      getCategoryList(),
      getTagList()
    ])
    
    if (categoryRes.code === 200) {
      categories.value = categoryRes.data || []
    }
    
    if (tagRes.code === 200) {
      allTags.value = tagRes.data || []
    }
    
    // 如果是编辑模式，获取文章详情
    if (isEdit.value) {
      await loadArticle(route.query.id)
    }
  } catch (error) {
    console.error('初始化数据失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 加载文章详情
const loadArticle = async (id) => {
  try {
    const response = await getArticleById(id)
    if (response.code === 200) {
      const data = response.data
      Object.assign(article, {
        id: data.id,
        title: data.title || '',
        slug: data.slug || '',
        summary: data.summary || '',
        content: data.content || '',
        categoryId: data.categoryId,
        tagIds: data.tagIds || [],
        thumbnail: data.thumbnail || '',
        metaTitle: data.meta_title || '',
        metaDescription: data.meta_description || '',
        metaKeywords: data.meta_keywords || '',
        isTop: !!data.isTop,
        isComment: !!data.isComment,
        status: data.status || 0
      })
      
      // 设置缩略图
      if (data.thumbnail) {
        fileList.value = [{
          uid: Date.now(),
          name: 'thumbnail',
          url: data.thumbnail
        }]
      }
      
      // 设置编辑器内容
      if (vditor.value && data.content) {
        vditor.value.setValue(data.content)
      }
      
      // 设置选中的标签
      if (data.tagIds && Array.isArray(data.tagIds)) {
        selectedTags.value = data.tagIds.map(tagId => {
          const tag = allTags.value.find(t => t.id === tagId)
          return tag ? tag.name : ''
        }).filter(name => name)
      }
    }
  } catch (error) {
    console.error('加载文章失败:', error)
    ElMessage.error('加载文章失败')
  }
}

// 生成slug
const generateSlug = () => {
  if (!article.slug && article.title) {
    article.slug = article.title
      .toLowerCase()
      .replace(/[^\w\s-]/g, '') // 移除特殊字符
      .replace(/[\s_-]+/g, '-') // 替换空格和下划线为连字符
      .replace(/^-+|-+$/g, '') // 移除首尾连字符
  }
}

// 处理返回
const handleBack = () => {
  router.push('/article')
}

// 预览文章
const handlePreview = () => {
  if (!article.content) {
    ElMessage.warning('请先输入文章内容')
    return
  }
  // Vditor 提供内置预览功能
  if (vditor.value) {
    const htmlContent = vditor.value.getHTML()
    const previewWindow = window.open('', '_blank')
    previewWindow.document.write(`
      <!DOCTYPE html>
      <html>
        <head>
          <title>文章预览</title>
          <meta charset="UTF-8">
          <style>
            body { 
              font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
              max-width: 800px;
              margin: 0 auto;
              padding: 20px;
              line-height: 1.6;
            }
          </style>
        </head>
        <body>
          <h1>${article.title || '无标题'}</h1>
          <div>${htmlContent}</div>
        </body>
      </html>
    `)
    previewWindow.document.close()
  }
}

// 缩略图上传处理
const beforeThumbnailUpload = (rawFile) => {
  if (rawFile.type.indexOf('image/') !== 0) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (rawFile.size / 1024 / 1024 > 5) {
    ElMessage.error('图片大小不能超过 5MB!')
    return false
  }
  return true
}

const handleThumbnailSuccess = (response, file) => {
  if (response.code === 200) {
    article.thumbnail = response.data.url
    ElMessage.success('缩略图上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

const handleUploadError = () => {
  ElMessage.error('上传失败')
}

const handlePictureCardPreview = (file) => {
  previewUrl.value = file.url
  previewVisible.value = true
}

const handleRemove = (file) => {
  article.thumbnail = ''
  fileList.value = []
}

// 发布文章
const handlePublish = async () => {
  if (!validateArticle()) return
  
  submitting.value = true
  try {
    const data = {
      title: article.title,
      slug: article.slug || generateSlugFromTitle(),
      summary: article.summary,
      content: article.content,
      categoryId: article.categoryId,
      tagIds: article.tagIds,
      add_tags: article.addTags, // 新增的标签名称数组
      thumbnail: article.thumbnail,
      meta_title: article.metaTitle,
      meta_description: article.metaDescription,
      meta_keywords: article.metaKeywords,
      isTop: article.isTop ? 1 : 0,
      isComment: article.isComment ? 1 : 0,
      status: 1 // 发布状态
    }
    
    let response
    if (isEdit.value) {
      response = await updateArticle(article.id, data)
    } else {
      response = await createArticle(data)
    }
    
    if (response.code === 200) {
      ElMessage.success(isEdit.value ? '文章更新成功' : '文章发布成功')
      router.push('/article')
    } else {
      ElMessage.error(response.message || '操作失败')
    }
  } catch (error) {
    console.error('发布文章失败:', error)
    ElMessage.error('操作失败，请重试')
  } finally {
    submitting.value = false
  }
}

// 保存草稿
const handleSaveDraft = async () => {
  if (!article.title.trim()) {
    ElMessage.warning('请输入文章标题')
    return
  }
  
  saving.value = true
  try {
    const data = {
      title: article.title,
      slug: article.slug || generateSlugFromTitle(),
      summary: article.summary,
      content: article.content,
      categoryId: article.categoryId,
      tagIds: article.tagIds,
      add_tags: article.addTags, // 新增的标签名称数组
      thumbnail: article.thumbnail,
      meta_title: article.metaTitle,
      meta_description: article.metaDescription,
      meta_keywords: article.metaKeywords,
      isTop: article.isTop ? 1 : 0,
      isComment: article.isComment ? 1 : 0,
      status: 0 // 草稿状态
    }
    
    let response
    if (isEdit.value) {
      response = await updateArticle(article.id, data)
    } else {
      response = await createArticle(data)
    }
    
    if (response.code === 200) {
      ElMessage.success('草稿保存成功')
      if (!isEdit.value && response.data.id) {
        article.id = response.data.id
        router.replace({ path: '/publish', query: { id: response.data.id } })
      }
    } else {
      ElMessage.error(response.message || '保存失败')
    }
  } catch (error) {
    console.error('保存草稿失败:', error)
    ElMessage.error('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

// 文章验证
const validateArticle = () => {
  if (!article.title.trim()) {
    ElMessage.warning('请输入文章标题')
    return false
  }
  
  if (!article.content.trim()) {
    ElMessage.warning('请输入文章内容')
    return false
  }
  
  if (!article.categoryId) {
    ElMessage.warning('请选择文章分类')
    return false
  }
  
  return true
}

const generateSlugFromTitle = () => {
  return article.title
    .toLowerCase()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

// 组件挂载时初始化数据
onMounted(async () => {
  await initData()
  // 延迟初始化编辑器，确保DOM已渲染
  nextTick(() => {
    initVditor()
  })
})

// 组件卸载时清理编辑器
onBeforeUnmount(() => {
  if (vditor.value) {
    vditor.value.destroy()
    vditor.value = null
  }
})
</script>

<style scoped>
.article-editor-page {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.editor-container {
  display: flex;
  gap: 20px;
  align-items: flex-start; /* 顶部对齐 */
}

.editor-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 标题输入框样式 */
.title-input-wrapper {
  background: #fff;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  padding: 16px;
}

.title-input {
  font-size: 18px;
  font-weight: 500;
}

.title-input :deep(.el-input__wrapper) {
  border: none;
  box-shadow: none;
  padding: 0;
  background: transparent;
}

.title-input :deep(.el-input__wrapper:focus-within) {
  box-shadow: none;
}

/* 编辑器包装器 */
.editor-wrapper {
  flex: 1;
  background: #fff;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  overflow: hidden;
}

.vditor-container {
  min-height: 700px;
}

.editor-sidebar {
  width: 280px;
  flex-shrink: 0;
}

/* 文章配置表单样式 */
.article-config-form {
  background: #fff;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  overflow: hidden;
}

.config-form-content {
  padding: 20px;
}

.form-item {
  margin-bottom: 10px;
}

.form-item:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 8px;
}

/* 标签显示样式 */
.selected-tags-display {
  margin-top: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tag-item-display {
  margin: 0;
}

/* 开关选项区域 */
.switch-section {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 4px;
  margin-bottom: 0 !important;
}

.switch-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.switch-item:last-child {
  margin-bottom: 0;
}

.switch-label {
  font-size: 14px;
  color: #606266;
}

/* 上传区域 */
.upload-section .upload-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 4px;
  border: 1px dashed #dcdfe6;
}

.upload-section {
  margin-bottom: 0 !important;
}

.thumbnail-upload :deep(.el-upload--picture-card) {
  width: 120px;
  height: 120px;
  border: 1px dashed #d9d9d9;
}

.thumbnail-upload :deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 120px;
  height: 120px;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  text-align: center;
}

.new-tag-indicator {
  font-size: 10px;
  opacity: 0.8;
}

/* 响应式调整 */
@media (max-width: 1200px) {
  .editor-container {
    flex-direction: column;
  }
  
  .editor-sidebar {
    width: 100%;
  }
  
  .header-actions {
    flex-wrap: wrap;
    gap: 8px;
  }
}

@media (max-width: 768px) {
  .article-editor-page {
    padding: 10px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-start;
  }
  
  .vditor-container {
    min-height: 500px;
  }
  
  .editor-container {
    gap: 12px;
  }
  
  .config-form-content {
    padding: 16px;
  }
  
  .form-item {
    margin-bottom: 16px;
  }
}
</style>