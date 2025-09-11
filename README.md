# Matuto Blog

🚀 **现代化个人博客系统** - 基于 Go + Gin + Tailwind CSS 的轻量级博客平台

## 📋 项目概述

Matuto Blog 是一个功能完整的个人博客系统，采用 Go 语言开发，支持 Markdown 文章写作，提供美观的响应式界面，满足个人博客和内容创作的需求。

### ✨ 核心特性

- 📝 **Markdown 支持** - 完整的 Markdown 渲染，支持代码高亮、表格、任务列表
- 🎨 **响应式设计** - 基于 Tailwind CSS 的现代化 UI，完美适配各种设备
- 🔍 **搜索功能** - 全文搜索，快速找到相关文章
- 📂 **分类管理** - 灵活的文章分类和标签系统
- 💖 **互动功能** - 文章点赞、浏览量统计、评论系统
- 🏷️ **标签系统** - 多标签支持，便于内容组织
- 📊 **排序功能** - 支持按时间、热度等多种方式排序
- 🔧 **管理后台** - 完整的后台管理系统，支持文章发布、编辑

## 🏗️ 技术架构

### 后端技术栈
- **语言**: Go 1.23+
- **Web框架**: Gin
- **数据库**: MySQL 5.7+ / SQLite
- **ORM**: GORM v1.30.0
- **模板引擎**: Go HTML Template
- **Markdown渲染**: goldmark
- **配置管理**: Viper
- **日志**: Logrus
- **认证**: JWT

### 前端技术栈
- **CSS框架**: Tailwind CSS
- **图标**: Font Awesome
- **JavaScript**: 原生 JavaScript
- **响应式**: Mobile-First 设计
- **模板**: Go Template

## 📁 项目结构

```
MatutoBlog/
├── cmd/                      # 应用入口
├── config/                   # 配置文件
├── internal/                 # 内部包
│   ├── api/                 # API层
│   │   ├── controllers/     # 控制器
│   │   ├── middlewares/     # 中间件
│   │   └── router/          # 路由配置
│   ├── database/            # 数据库层
│   └── models/              # 数据模型
├── pkg/                     # 公共包
│   ├── common/              # 通用工具
│   ├── logger/              # 日志工具
│   ├── storage/             # 存储工具
│   └── utils/               # 工具函数
├── web/                     # Web资源
│   ├── static/              # 静态文件
│   ├── templates/           # 模板文件
│   │   └── default/         # 默认主题
│   │       ├── components/  # 组件模板
│   │       ├── index.html   # 首页
│   │       ├── article.html # 文章详情
│   │       └── category.html# 分类页面
│   └── uploads/             # 上传文件
├── docs/                    # 文档
└── scripts/                 # 脚本文件
```

## 🚀 快速开始

### 环境要求

- **Go**: 1.23+
- **MySQL**: 5.7+ 或 SQLite 3
- **现代浏览器**: Chrome, Firefox, Safari, Edge

### 1. 克隆项目

```bash
git clone <repository-url>
cd MatutoBlog
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置数据库

#### 使用 MySQL
```sql
CREATE DATABASE matuto_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 使用 SQLite（默认）
项目默认使用 SQLite，无需额外配置。

### 4. 配置文件

根据需要修改配置文件或设置环境变量：

```bash
# 服务器配置
export SERVER_PORT=8080
export SERVER_MODE=debug

# 数据库配置（MySQL）
export DATABASE_HOST=localhost
export DATABASE_PORT=3306
export DATABASE_USERNAME=root
export DATABASE_PASSWORD=password
export DATABASE_DBNAME=matuto_blog

# 或使用 SQLite（默认）
export DATABASE_TYPE=sqlite
export DATABASE_PATH=./blog.db
```

### 5. 启动服务

```bash
# 开发模式
go run main.go

# 或构建后运行
go build -o blog.exe .
./blog.exe
```

### 6. 访问应用

- 前台博客: http://localhost:8080
- 管理后台: http://localhost:8080/admin（需要登录）

## 📚 功能使用

### 1. 文章管理

#### 创建文章
- 支持 Markdown 语法
- 文章分类和标签
- 文章封面图片
- 发布时间设置
- 置顶功能

#### 文章特性
- **Markdown 渲染**: 支持 GitHub Flavored Markdown
- **代码高亮**: 预设代码块样式
- **表格支持**: 完整的表格渲染
- **任务列表**: 支持 Todo 列表
- **数学公式**: 支持基础数学符号

### 2. 分类系统

#### 分类管理
- 创建、编辑、删除分类
- 分类描述和封面
- 文章数量统计
- 分类页面展示

#### 标签系统
- 多标签关联
- 标签云展示
- 按标签筛选文章

### 3. 搜索功能

- 全文搜索支持
- 标题和内容搜索
- 搜索结果高亮
- 搜索历史记录

### 4. 互动功能

- **点赞系统**: 文章点赞统计
- **浏览统计**: 实时浏览量记录
- **评论系统**: 支持文章评论（静态展示）
- **社交分享**: 内置社交媒体分享

## 🔧 开发指南

### 1. 添加新页面

1. 在 `web/templates/default/` 下创建模板文件
2. 在 `internal/api/controllers/` 下添加控制器方法
3. 在 `internal/api/router/routes.go` 中添加路由

### 2. 自定义主题

1. 复制 `web/templates/default/` 目录
2. 重命名为新主题名称
3. 修改模板文件和样式
4. 在配置中切换主题

### 3. 扩展功能

#### 添加新的模板函数
在 `pkg/utils/template.go` 的 `GenTemplateFuncMap()` 中添加：

```go
"customFunc": func(input string) string {
    // 自定义逻辑
    return output
},
```

#### 添加新的中间件
在 `internal/api/middlewares/` 下创建中间件文件。

### 4. 数据库模型

在 `internal/models/` 下定义新的数据模型：

```go
type NewModel struct {
    models.BaseModel
    Name        string `json:"name" gorm:"size:100;not null"`
    Description string `json:"description" gorm:"type:text"`
    Status      int    `json:"status" gorm:"default:1"`
}
```

## 📊 API 接口

### 前台接口

- `GET /` - 首页文章列表
- `GET /article/:id` - 文章详情
- `GET /category/:id` - 分类页面
- `GET /categories` - 分类列表
- `GET /tag/:id` - 标签页面
- `GET /search` - 搜索页面

### 管理接口

- `POST /api/login` - 管理员登录
- `GET /api/articles/page` - 文章分页列表
- `POST /api/articles/publish` - 发布文章
- `PUT /api/articles/update` - 更新文章
- `DELETE /api/articles/:id` - 删除文章
- `GET /api/categories/page` - 分类管理
- `POST /api/categories` - 创建分类
- `GET /api/tags/page` - 标签管理

## 🚀 部署指南

### 1. 生产构建

```bash
# 构建可执行文件
go build -ldflags="-w -s" -o blog .

# 压缩体积（可选）
upx --brute blog
```

### 2. 系统服务配置

创建 systemd 服务文件 `/etc/systemd/system/matuto-blog.service`：

```ini
[Unit]
Description=Matuto Blog Service
After=network.target

[Service]
Type=simple
User=blog
WorkingDirectory=/opt/matuto-blog
ExecStart=/opt/matuto-blog/blog
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl enable matuto-blog
sudo systemctl start matuto-blog
```

### 3. Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /static/ {
        alias /opt/matuto-blog/web/static/;
        expires 30d;
        add_header Cache-Control "public, no-transform";
    }

    location /uploads/ {
        alias /opt/matuto-blog/web/uploads/;
        expires 30d;
        add_header Cache-Control "public, no-transform";
    }
}
```

### 4. Docker 部署

创建 `Dockerfile`：

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o blog .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/blog .
COPY --from=builder /app/web ./web
COPY --from=builder /app/config ./config
EXPOSE 8080
CMD ["./blog"]
```

构建和运行：
```bash
docker build -t matuto-blog .
docker run -d -p 8080:8080 --name blog matuto-blog
```

## 📈 性能优化

### 1. 数据库优化
- 添加适当的索引
- 使用连接池
- 查询优化

### 2. 缓存策略
- 静态文件缓存
- 页面缓存
- 数据缓存

### 3. 前端优化
- 图片压缩和懒加载
- CSS/JS 压缩
- CDN 使用

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 开源协议

本项目基于 MIT 协议开源 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 常见问题

### Q: 如何修改博客主题？
A: 复制 `web/templates/default` 目录，重命名后修改模板文件，在配置中指定新主题名。

### Q: 如何添加自定义页面？
A: 在模板目录下创建新的 HTML 文件，在控制器中添加处理方法，在路由中注册新路由。

### Q: 数据库迁移怎么处理？
A: 项目启动时会自动执行数据库迁移，如需手动执行可使用 GORM 的 AutoMigrate 功能。

### Q: 如何备份数据？
A: 定期备份数据库文件（SQLite）或使用 mysqldump（MySQL）备份数据库。

## 📞 联系方式

- 项目主页: [Matuto Blog](https://github.com/your-username/MatutoBlog)
- 问题反馈: [Issues](https://github.com/your-username/MatutoBlog/issues)
- 邮箱: [your-email@example.com]

---

**Matuto Blog** - 让写作更简单，让分享更美好！ ✨