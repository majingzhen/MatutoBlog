<template>
  <div class="article-editor-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <el-page-header @back="handleBack" :content="isEdit ? '编辑文章' : '新建文章'" />
      <div class="header-actions">
        <el-button @click="handlePreview" type="info" plain>
          <el-icon><View /></el-icon>
          预览
        </el-button>
      </div>
    </div>

    <div class="editor-container">
      <!-- 左侧编辑区域 -->
      <div class="editor-main">
        <!-- 标题输入框 -->
        <el-card class="title-input-card" shadow="never">
          <el-input
            v-model="article.title"
            placeholder="请输入文章标题"
            size="large"
            :maxlength="200"
            show-word-limit
            class="title-input"
            @blur="generateSlug"
          />
        </el-card>

        <!-- Markdown编辑器 -->
        <el-card class="editor-card" shadow="never">
          <el-tabs v-model="activeTab" type="card" class="editor-tabs">
            <el-tab-pane label="编辑" name="edit">
              <div class="markdown-editor-container">
                <div class="editor-toolbar">
                  <el-button-group>
                    <el-button size="small" @click="insertMarkdown('**', '**')" title="加粗">
                      <el-icon><bold /></el-icon>
                    </el-button>
                    <el-button size="small" @click="insertMarkdown('*', '*')" title="斜体">
                      <el-icon><italic /></el-icon>
                    </el-button>
                    <el-button size="small" @click="insertMarkdown('# ', '')" title="标题">
                      H
                    </el-button>
                    <el-button size="small" @click="insertMarkdown('`', '`')" title="代码">
                      <el-icon><code /></el-icon>
                    </el-button>
                    <el-button size="small" @click="insertMarkdown('[', '](url)')" title="链接">
                      <el-icon><link /></el-icon>
                    </el-button>
                    <el-button size="small" @click="insertMarkdown('![', '](image-url)')" title="图片">
                      <el-icon><picture /></el-icon>
                    </el-button>
                    <el-button size="small" @click="insertMarkdown('- ', '')" title="列表">
                      <el-icon><list /></el-icon>
                    </el-button>
                    <el-button size="small" @click="insertMarkdown('> ', '')" title="引用">
                      <el-icon><chat-quote-fill /></el-icon>
                    </el-button>
                  </el-button-group>
                  
                  <el-upload
                    :action="uploadUrl"
                    :headers="uploadHeaders"
                    :show-file-list="false"
                    :before-upload="beforeImageUpload"
                    :on-success="handleImageSuccess"
                    accept="image/*"
                    style="display: inline-block; margin-left: 10px;"
                  >
                    <el-button size="small" type="primary">
                      <el-icon><upload /></el-icon>
                      上传图片
                    </el-button>
                  </el-upload>
                </div>
                
                <textarea
                  ref="markdownTextarea"
                  v-model="article.content"
                  class="markdown-editor"
                  placeholder="开始创作你的文章..."
                  @input="handleContentChange"
                  @scroll="syncScroll"
                />
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="预览" name="preview">
              <div 
                ref="previewContainer"
                class="markdown-preview"
                v-html="renderedContent"
                @scroll="syncScroll"
              />
            </el-tab-pane>
            
            <el-tab-pane label="分屏" name="split">
              <div class="split-editor">
                <div class="split-edit">
                  <div class="editor-toolbar">
                    <el-button-group>
                      <el-button size="small" @click="insertMarkdown('**', '**')" title="加粗">
                        <el-icon><bold /></el-icon>
                      </el-button>
                      <el-button size="small" @click="insertMarkdown('*', '*')" title="斜体">
                        <el-icon><italic /></el-icon>
                      </el-button>
                      <el-button size="small" @click="insertMarkdown('# ', '')" title="标题">
                        H
                      </el-button>
                      <el-button size="small" @click="insertMarkdown('`', '`')" title="代码">
                        <el-icon><code /></el-icon>
                      </el-button>
                      <el-button size="small" @click="insertMarkdown('[', '](url)')" title="链接">
                        <el-icon><link /></el-icon>
                      </el-button>
                      <el-button size="small" @click="insertMarkdown('![', '](image-url)')" title="图片">
                        <el-icon><picture /></el-icon>
                      </el-button>
                    </el-button-group>
                    
                    <el-upload
                      :action="uploadUrl"
                      :headers="uploadHeaders"
                      :show-file-list="false"
                      :before-upload="beforeImageUpload"
                      :on-success="handleImageSuccess"
                      accept="image/*"
                      style="display: inline-block; margin-left: 10px;"
                    >
                      <el-button size="small" type="primary">
                        <el-icon><upload /></el-icon>
                        上传图片
                      </el-button>
                    </el-upload>
                  </div>
                  <textarea
                    v-model="article.content"
                    class="split-textarea"
                    placeholder="开始创作你的文章..."
                    @input="handleContentChange"
                  />
                </div>
                <div class="split-preview">
                  <div class="preview-header">预览</div>
                  <div 
                    class="split-preview-content"
                    v-html="renderedContent"
                  />
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </div>

      <!-- 右侧配置区域 -->
      <div class="editor-sidebar">
        <!-- 文章摘要 -->
        <el-card class="sidebar-card" shadow="never">
          <template #header>
            <div class="sidebar-title">
              <el-icon><Document /></el-icon>
              <span>文章摘要</span>
            </div>
          </template>
          <el-input
            v-model="article.summary"
            type="textarea"
            :rows="3"
            placeholder="请输入文章摘要..."
            :maxlength="500"
            show-word-limit
          />
        </el-card>

        <!-- 分类选择 -->
        <el-card class="sidebar-card" shadow="never">
          <template #header>
            <div class="sidebar-title">
              <el-icon><Folder /></el-icon>
              <span>文章分类</span>
            </div>
          </template>
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
        </el-card>

        <!-- 标签选择 -->
        <el-card class="sidebar-card" shadow="never">
          <template #header>
            <div class="sidebar-title">
              <el-icon><Tags /></el-icon>
              <span>文章标签</span>
            </div>
          </template>
          <el-select
            v-model="article.tagIds"
            multiple
            placeholder="请选择标签"
            style="width: 100%;"
            filterable
            allow-create
            :reserve-keyword="false"
            default-first-option
          >
            <el-option
              v-for="tag in allTags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
          <p class="tags-hint">提示：可以创建新标签</p>
        </el-card>

        <!-- 缩略图 -->
        <el-card class="sidebar-card" shadow="never">
          <template #header>
            <div class="sidebar-title">
              <el-icon><Picture /></el-icon>
              <span>缩略图</span>
            </div>
          </template>
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
        </el-card>

        <!-- SEO设置 -->
        <el-card class="sidebar-card" shadow="never">
          <template #header>
            <div class="sidebar-title">
              <el-icon><Search /></el-icon>
              <span>SEO设置</span>
            </div>
          </template>
          <el-form label-position="top">
            <el-form-item label="SEO标题">
              <el-input
                v-model="article.metaTitle"
                placeholder="留空将使用文章标题"
                :maxlength="60"
                show-word-limit
              />
            </el-form-item>
            <el-form-item label="SEO描述">
              <el-input
                v-model="article.metaDescription"
                placeholder="留空将使用文章摘要"
                type="textarea"
                :rows="3"
                :maxlength="160"
                show-word-limit
              />
            </el-form-item>
            <el-form-item label="关键词">
              <el-input
                v-model="article.metaKeywords"
                placeholder="用逗号分隔关键词"
                :maxlength="200"
                show-word-limit
              />
            </el-form-item>
            <el-form-item label="文章别名(Slug)">
              <el-input
                v-model="article.slug"
                placeholder="文章URL别名"
                :maxlength="100"
              />
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 文章设置 -->
        <el-card class="sidebar-card" shadow="never">
          <template #header>
            <div class="sidebar-title">
              <el-icon><Setting /></el-icon>
              <span>文章设置</span>
            </div>
          </template>
          <el-form label-position="top">
            <el-form-item>
              <el-checkbox v-model="article.isTop">置顶文章</el-checkbox>
            </el-form-item>
            <el-form-item>
              <el-checkbox v-model="article.isComment">允许评论</el-checkbox>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 操作按钮 -->
        <div class="action-buttons">
          <el-button 
            type="primary" 
            size="large"
            class="publish-btn"
            @click="handlePublish"
            :loading="submitting"
          >
            <el-icon><Upload /></el-icon>
            {{ isEdit ? '更新文章' : '发布文章' }}
          </el-button>
          <el-button 
            size="large"
            class="save-btn"
            @click="handleSaveDraft"
            :loading="saving"
          >
            <el-icon><Save /></el-icon>
            保存草稿
          </el-button>
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
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  View, Document, Folder, Picture, Search, Setting,
  Plus, Delete, ZoomIn, Upload,
   Link, List
} from '@element-plus/icons-vue'

// 导入markdown渲染库
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

// 导入API
import { 
  createArticle, updateArticle, getArticleById,
  getCategoryList, getTagList
} from '@/api/article.js'

// 初始化marked
marked.setOptions({
  highlight: function(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      return hljs.highlight(code, { language: lang }).value
    }
    return hljs.highlightAuto(code).value
  },
  breaks: true,
  gfm: true
})

const router = useRouter()
const route = useRoute()

// 判断是否为编辑模式
const isEdit = computed(() => !!route.query.id)

// 状态控制
const submitting = ref(false)
const saving = ref(false)
const loading = ref(false)
const activeTab = ref('edit')

// 编辑器相关
const markdownTextarea = ref(null)
const previewContainer = ref(null)
const renderedContent = ref('')

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
  thumbnail: '',
  metaTitle: '',
  metaDescription: '',
  metaKeywords: '',
  isTop: false,
  isComment: true,
  status: 0 // 0: 草稿, 1: 发布
})

// 分类和标签数据
const categories = ref([])
const allTags = ref([])

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
  renderedContent.value = marked.parse(article.content || '')
}

// 处理内容变化
const handleContentChange = () => {
  renderMarkdown()
  article.parseContent = renderedContent.value
}

// 插入Markdown语法
const insertMarkdown = (before, after) => {
  const textarea = markdownTextarea.value
  if (!textarea) return
  
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = textarea.value.substring(start, end)
  
  const newText = before + selectedText + after
  const newContent = 
    textarea.value.substring(0, start) + 
    newText + 
    textarea.value.substring(end)
  
  article.content = newContent
  handleContentChange()
  
  // 恢复光标位置
  nextTick(() => {
    textarea.focus()
    const newCursorPos = start + before.length + selectedText.length
    textarea.setSelectionRange(newCursorPos, newCursorPos)
  })
}

// 同步滚动
const syncScroll = (e) => {
  const source = e.target
  const target = activeTab.value === 'edit' ? previewContainer.value : markdownTextarea.value
  if (target && source.scrollHeight > source.clientHeight) {
    const scrollRatio = source.scrollTop / (source.scrollHeight - source.clientHeight)
    target.scrollTop = scrollRatio * (target.scrollHeight - target.clientHeight)
  }
}

// 图片上传处理
const beforeImageUpload = (rawFile) => {
  if (rawFile.type.indexOf('image/') !== 0) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (rawFile.size / 1024 / 1024 > 10) {
    ElMessage.error('图片大小不能超过 10MB!')
    return false
  }
  return true
}

const handleImageSuccess = (response, file) => {
  if (response.code === 200) {
    const imageUrl = response.data.url
    const markdownImage = `![${file.name}](${imageUrl})`
    insertMarkdown('', markdownImage)
    ElMessage.success('图片上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
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
        categoryId: data.category_id,
        tagIds: data.tag_ids || [],
        thumbnail: data.thumbnail || '',
        metaTitle: data.meta_title || '',
        metaDescription: data.meta_description || '',
        metaKeywords: data.meta_keywords || '',
        isTop: !!data.is_top,
        isComment: !!data.is_comment,
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
      
      // 渲染内容
      handleContentChange()
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
  window.open(`data:text/html,${encodeURIComponent(renderedContent.value)}`, '_blank')
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
      category_id: article.categoryId,
      tag_ids: article.tagIds,
      thumbnail: article.thumbnail,
      meta_title: article.metaTitle,
      meta_description: article.metaDescription,
      meta_keywords: article.metaKeywords,
      is_top: article.isTop ? 1 : 0,
      is_comment: article.isComment ? 1 : 0,
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
      category_id: article.categoryId,
      tag_ids: article.tagIds,
      thumbnail: article.thumbnail,
      meta_title: article.metaTitle,
      meta_description: article.metaDescription,
      meta_keywords: article.metaKeywords,
      is_top: article.isTop ? 1 : 0,
      is_comment: article.isComment ? 1 : 0,
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
onMounted(() => {
  initData()
  renderMarkdown()
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
  gap: 10px;
}

.editor-container {
  display: flex;
  gap: 20px;
}

.editor-main {
  flex: 1;
  min-width: 0;
}

.title-input-card {
  margin-bottom: 15px;
}

.title-input {
  font-size: 22px;
  font-weight: 500;
}

.title-input :deep(.el-input__wrapper) {
  border: none;
  box-shadow: none;
  padding: 20px 0;
}

.title-input :deep(.el-input__wrapper:focus-within) {
  box-shadow: none;
}

.editor-card {
  overflow: hidden;
}

.editor-card :deep(.el-card__body) {
  padding: 0;
}

.editor-tabs :deep(.el-tabs__content) {
  padding: 0;
}

.markdown-editor-container {
  position: relative;
}

.editor-toolbar {
  padding: 10px 15px;
  border-bottom: 1px solid #e5e7eb;
  background-color: #f8fafc;
  display: flex;
  align-items: center;
  gap: 10px;
}

.markdown-editor {
  width: 100%;
  height: 500px;
  padding: 20px;
  border: none;
  resize: none;
  font-family: 'Consolas', 'Monaco', 'Menlo', monospace;
  font-size: 14px;
  line-height: 1.6;
  outline: none;
}

.markdown-preview {
  padding: 20px;
  height: 500px;
  overflow-y: auto;
  background-color: #fff;
}

.split-editor {
  display: flex;
  height: 500px;
}

.split-edit {
  flex: 1;
  border-right: 1px solid #e5e7eb;
}

.split-textarea {
  width: 100%;
  height: calc(100% - 50px);
  padding: 15px;
  border: none;
  resize: none;
  font-family: 'Consolas', 'Monaco', 'Menlo', monospace;
  font-size: 14px;
  line-height: 1.6;
  outline: none;
}

.split-preview {
  flex: 1;
}

.preview-header {
  padding: 10px 15px;
  background-color: #f8fafc;
  border-bottom: 1px solid #e5e7eb;
  font-weight: 500;
  font-size: 14px;
}

.split-preview-content {
  padding: 15px;
  height: calc(100% - 50px);
  overflow-y: auto;
}

.editor-sidebar {
  width: 350px;
  flex-shrink: 0;
}

.sidebar-card {
  margin-bottom: 20px;
}

.sidebar-title {
  display: flex;
  align-items: center;
  font-weight: 500;
  color: #1d2129;
}

.sidebar-title span {
  margin-left: 5px;
}

.tags-hint {
  margin-top: 10px;
  font-size: 12px;
  color: #86909c;
  margin-bottom: 0;
}

.action-buttons {
  position: sticky;
  top: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.publish-btn, .save-btn {
  width: 100%;
}

/* Markdown预览样式 */
.markdown-preview :deep(h1),
.split-preview-content :deep(h1) {
  font-size: 2em;
  margin: 0.67em 0;
  border-bottom: 1px solid #eaecef;
  padding-bottom: 0.3em;
}

.markdown-preview :deep(h2),
.split-preview-content :deep(h2) {
  font-size: 1.5em;
  margin: 0.83em 0;
  border-bottom: 1px solid #eaecef;
  padding-bottom: 0.3em;
}

.markdown-preview :deep(h3),
.split-preview-content :deep(h3) {
  font-size: 1.17em;
  margin: 1em 0;
}

.markdown-preview :deep(p),
.split-preview-content :deep(p) {
  margin: 1em 0;
  line-height: 1.6;
}

.markdown-preview :deep(code),
.split-preview-content :deep(code) {
  background-color: #f1f3f4;
  padding: 0.2em 0.4em;
  border-radius: 3px;
  font-family: 'Consolas', 'Monaco', 'Menlo', monospace;
}

.markdown-preview :deep(pre),
.split-preview-content :deep(pre) {
  background-color: #f6f8fa;
  padding: 1em;
  border-radius: 6px;
  overflow-x: auto;
  margin: 1em 0;
}

.markdown-preview :deep(pre code),
.split-preview-content :deep(pre code) {
  background-color: transparent;
  padding: 0;
}

.markdown-preview :deep(blockquote),
.split-preview-content :deep(blockquote) {
  border-left: 4px solid #dfe2e5;
  padding-left: 1em;
  margin: 1em 0;
  color: #6a737d;
}

.markdown-preview :deep(ul),
.markdown-preview :deep(ol),
.split-preview-content :deep(ul),
.split-preview-content :deep(ol) {
  margin: 1em 0;
  padding-left: 2em;
}

.markdown-preview :deep(table),
.split-preview-content :deep(table) {
  border-collapse: collapse;
  margin: 1em 0;
  width: 100%;
}

.markdown-preview :deep(th),
.markdown-preview :deep(td),
.split-preview-content :deep(th),
.split-preview-content :deep(td) {
  border: 1px solid #d0d7de;
  padding: 0.5em;
}

.markdown-preview :deep(th),
.split-preview-content :deep(th) {
  background-color: #f6f8fa;
  font-weight: 600;
}

/* 上传组件样式 */
:deep(.el-upload--picture-card) {
  width: 100px;
  height: 100px;
}

:deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 100px;
  height: 100px;
}

/* 响应式调整 */
@media (max-width: 1200px) {
  .editor-container {
    flex-direction: column;
  }
  
  .editor-sidebar {
    width: 100%;
  }
  
  .action-buttons {
    position: static;
    flex-direction: row;
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
  
  .action-buttons {
    flex-direction: column;
  }
  
  .split-editor {
    flex-direction: column;
  }
  
  .split-edit {
    border-right: none;
    border-bottom: 1px solid #e5e7eb;
  }
}
</style>