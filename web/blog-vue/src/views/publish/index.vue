<template>
  <div class="article-editor-page">
    <!-- 页面标题 -->
    <el-page-header 
      @back="handleBack"
      content="文章编辑"
    />

    <!-- 编辑器切换 -->
    <el-card class="editor-switch-card" :bordered="false">
      <div class="editor-switch">
        <span class="switch-label">编辑器类型：</span>
        <el-radio-group v-model="editorType" @change="handleEditorChange">
          <el-radio-button label="normal">普通编辑器</el-radio-button>
          <el-radio-button label="ai">AiEditor</el-radio-button>
        </el-radio-group>
      </div>
    </el-card>

    <div class="editor-container">
      <!-- 左侧编辑区域 -->
      <div class="editor-main">
        <!-- 标题输入框 -->
        <el-card class="title-input-card" :bordered="false">
          <el-input
            v-model="article.title"
            placeholder="请输入文章标题"
            size="large"
            :maxlength="100"
            show-word-limit
            class="title-input"
          />
        </el-card>

        <!-- 普通富文本编辑器 -->
        <el-card v-if="editorType === 'normal'" class="editor-card">
          <el-tabs v-model="normalEditorTab" type="card">
            <el-tab-pane label="编辑视图">
              <textarea
                v-model="article.content"
                class="markdown-editor"
                placeholder="请使用Markdown格式编写文章内容..."
              ></textarea>
            </el-tab-pane>
            <el-tab-pane label="预览">
              <div 
                class="markdown-preview"
                v-html="renderMarkdown(article.content)"
              ></div>
            </el-tab-pane>
          </el-tabs>
        </el-card>

        <!-- AiEditor -->
        <el-card v-if="editorType === 'ai'" class="editor-card">
          <div class="ai-editor-header">
            <el-button 
              type="text" 
              size="small"
              @click="handleAiGenerate"
            >
              <el-icon><MagicStick /></el-icon>
              AI辅助生成
            </el-button>
            <el-button 
              type="text" 
              size="small"
              @click="handleAiOptimize"
            >
              <el-icon><Brush /></el-icon>
              AI优化内容
            </el-button>
          </div>
          <div class="ai-editor-content">
            <textarea
              v-model="article.content"
              class="ai-markdown-editor"
              placeholder="请输入内容，或使用AI辅助生成..."
            ></textarea>
          </div>
          <div class="ai-editor-footer">
            <el-progress 
              v-if="aiProcessing" 
              percentage="50" 
              stroke-width="2" 
              style="width: 200px;"
            />
            <span v-if="aiMessage" class="ai-message">{{ aiMessage }}</span>
          </div>
        </el-card>
      </div>

      <!-- 右侧配置区域 -->
      <div class="editor-sidebar">
        <!-- 分类选择 -->
        <el-card class="sidebar-card" :bordered="false">
          <div class="sidebar-title">
            <el-icon><Folder /></el-icon>
            <span>文章分类</span>
          </div>
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
          <el-button 
            type="text" 
            size="small" 
            class="add-category"
            @click="handleAddCategory"
          >
            + 添加新分类
          </el-button>
        </el-card>

        <!-- 标签选择 -->
        <el-card class="sidebar-card" :bordered="false">
          <div class="sidebar-title">
            <el-icon><Tags /></el-icon>
            <span>文章标签</span>
          </div>
          <el-select
            v-model="article.tags"
            multiple
            placeholder="请选择或输入标签"
            style="width: 100%;"
            filterable
            allow-create
          >
            <el-option
              v-for="tag in allTags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.name"
            />
          </el-select>
          <p class="tags-hint">提示：输入后按回车添加自定义标签</p>
        </el-card>

        <!-- SEO设置 -->
        <el-card class="sidebar-card" :bordered="false">
          <div class="sidebar-title">
            <el-icon><Search /></el-icon>
            <span>SEO设置</span>
          </div>
          <el-form>
            <el-form-item label="SEO标题" size="small">
              <el-input
                v-model="article.seoTitle"
                placeholder="自定义SEO标题"
                :maxlength="60"
                show-word-limit
              />
            </el-form-item>
            <el-form-item label="SEO描述" size="small">
              <el-input
                v-model="article.seoDescription"
                placeholder="自定义SEO描述"
                type="textarea"
                :rows="3"
                :maxlength="160"
                show-word-limit
              />
            </el-form-item>
            <el-form-item label="关键词" size="small">
              <el-input
                v-model="article.seoKeywords"
                placeholder="用逗号分隔关键词"
                :maxlength="100"
                show-word-limit
              />
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
            发布文章
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
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElDialog, ElInput, ElButton } from 'element-plus'
import { 
  Folder, Tags, Search, Upload, Save, 
  MagicStick, Brush, ArrowLeft 
} from '@element-plus/icons-vue'

// 引入markdown渲染库（实际项目中需要安装：npm install marked highlight.js）
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

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

// 编辑器类型
const editorType = ref('normal')
const normalEditorTab = ref('edit')

// 状态控制
const submitting = ref(false)
const saving = ref(false)
const aiProcessing = ref(false)
const aiMessage = ref('')

// 文章数据
const article = reactive({
  id: '', // 新增时为空，编辑时有值
  title: '',
  content: '',
  categoryId: '',
  tags: [],
  seoTitle: '',
  seoDescription: '',
  seoKeywords: ''
})

// 分类数据
const categories = ref([
  { id: 1, name: '前端开发' },
  { id: 2, name: '后端开发' },
  { id: 3, name: '数据库' },
  { id: 4, name: '服务器' },
  { id: 5, name: '人工智能' }
])

// 标签数据
const allTags = ref([
  { id: 1, name: 'JavaScript' },
  { id: 2, name: 'Vue' },
  { id: 3, name: 'React' },
  { id: 4, name: 'Node.js' },
  { id: 5, name: 'HTML' },
  { id: 6, name: 'CSS' },
  { id: 7, name: 'Go' }
])

// 新分类对话框
const categoryDialog = ref(false)
const newCategory = ref('')

// 渲染markdown
const renderMarkdown = (content) => {
  return marked.parse(content || '')
}

// 处理返回
const handleBack = () => {
  router.back()
}

// 切换编辑器
const handleEditorChange = (type) => {
  console.log('切换到', type, '编辑器')
  // 实际项目中可以在这里处理两种编辑器内容的转换
}

// AI生成内容
const handleAiGenerate = () => {
  if (!article.title) {
    ElMessage.warning('请先输入文章标题')
    return
  }
  
  aiProcessing.value = true
  aiMessage.value = 'AI正在生成内容...'
  
  // 模拟AI生成
  setTimeout(() => {
    article.content += `\n## 关于${article.title}的介绍\n\n`
    article.content += '本文将详细介绍相关内容，包括以下几个方面：\n\n'
    article.content += '- 基本概念和原理\n'
    article.content += '- 实际应用场景\n'
    article.content += '- 常见问题及解决方案\n'
    article.content += '- 未来发展趋势\n\n'
    
    aiProcessing.value = false
    aiMessage.value = '内容生成完成，可以继续编辑'
    
    setTimeout(() => {
      aiMessage.value = ''
    }, 3000)
  }, 2000)
}

// AI优化内容
const handleAiOptimize = () => {
  if (!article.content) {
    ElMessage.warning('请先输入内容')
    return
  }
  
  aiProcessing.value = true
  aiMessage.value = 'AI正在优化内容...'
  
  // 模拟AI优化
  setTimeout(() => {
    aiProcessing.value = false
    aiMessage.value = '内容优化完成'
    
    setTimeout(() => {
      aiMessage.value = ''
    }, 2000)
  }, 1500)
}

// 添加新分类
const handleAddCategory = () => {
  categoryDialog.value = true
  newCategory.value = ''
}

// 确认添加分类
const confirmAddCategory = () => {
  if (!newCategory.value.trim()) {
    ElMessage.warning('分类名称不能为空')
    return
  }
  
  // 模拟添加分类
  const newId = categories.value.length + 1
  categories.value.push({
    id: newId,
    name: newCategory.value.trim()
  })
  
  // 自动选中新添加的分类
  article.categoryId = newId
  
  categoryDialog.value = false
  ElMessage.success('分类添加成功')
}

// 发布文章
const handlePublish = () => {
  // 简单验证
  if (!validateArticle()) return
  
  submitting.value = true
  
  // 模拟发布请求
  setTimeout(() => {
    ElMessage.success('文章发布成功')
    submitting.value = false
    router.push('/article/list')
  }, 1500)
}

// 保存草稿
const handleSaveDraft = () => {
  // 草稿可以不验证完整信息，但至少需要标题
  if (!article.title.trim()) {
    ElMessage.warning('请输入文章标题')
    return
  }
  
  saving.value = true
  
  // 模拟保存请求
  setTimeout(() => {
    ElMessage.success('草稿保存成功')
    saving.value = false
  }, 1000)
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
</script>

<style scoped>
.article-editor-page {
  padding: 20px;
}

.editor-switch-card {
  margin-bottom: 15px;
  background-color: #f5f7fa;
}

.editor-switch {
  display: flex;
  align-items: center;
  padding: 10px 0;
}

.switch-label {
  margin-right: 15px;
  color: #606266;
}

.editor-container {
  display: flex;
  gap: 20px;
}

.editor-main {
  flex: 1;
}

.title-input-card {
  margin-bottom: 15px;
}

.title-input {
  font-size: 20px;
  font-weight: 500;
}

.title-input .el-input__wrapper {
  border: none;
  box-shadow: none;
  padding: 15px 0;
}

.title-input .el-input__wrapper:focus-within {
  box-shadow: none;
}

.editor-card {
  min-height: 500px;
}

.markdown-editor {
  width: 100%;
  height: 450px;
  padding: 15px;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  resize: vertical;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  line-height: 1.6;
}

.markdown-editor:focus {
  outline: none;
  border-color: #409eff;
}

.markdown-preview {
  padding: 15px;
  min-height: 450px;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  overflow-y: auto;
}

.markdown-preview h1,
.markdown-preview h2,
.markdown-preview h3 {
  margin: 1.5rem 0 1rem;
}

.markdown-preview p {
  margin-bottom: 1rem;
  line-height: 1.8;
}

.markdown-preview ul,
.markdown-preview ol {
  margin-bottom: 1rem;
  padding-left: 2rem;
}

.markdown-preview pre {
  background-color: #f8fafc;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  overflow-x: auto;
}

.markdown-preview code {
  font-family: 'Consolas', 'Monaco', monospace;
}

.ai-editor-header {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e5e7eb;
}

.ai-markdown-editor {
  width: 100%;
  height: 400px;
  padding: 15px;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  resize: vertical;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  line-height: 1.6;
}

.ai-markdown-editor:focus {
  outline: none;
  border-color: #409eff;
}

.ai-editor-footer {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #e5e7eb;
}

.ai-message {
  color: #409eff;
  font-size: 13px;
}

.editor-sidebar {
  width: 350px;
}

.sidebar-card {
  margin-bottom: 20px;
}

.sidebar-title {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  font-weight: 500;
  color: #1d2129;
}

.sidebar-title span {
  margin-left: 5px;
}

.add-category {
  color: #409eff;
  padding: 5px 0;
}

.tags-hint {
  margin-top: 10px;
  font-size: 12px;
  color: #86909c;
  margin-bottom: 0;
}

.el-form-item__label {
  font-size: 13px;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.publish-btn, .save-btn {
  width: 100%;
}

/* 响应式调整 */
@media (max-width: 992px) {
  .editor-container {
    flex-direction: column;
  }
  
  .editor-sidebar {
    width: 100%;
  }
}
</style>
