# 博客系统数据库设计文档

## 1. 概述

本文档描述了一个简单的个人博客系统的数据库设计。该系统只有一个用户（博主），分为管理端和前台展示端。前台所有操作不需要登录，但评论需要审核。

## 2. 设计原则

1. 简洁性：针对单用户博客系统进行简化设计
2. 完整性：包含博客系统所需的核心功能
3. 性能：为常用查询字段添加索引
4. 安全性：增强密码存储安全性

## 3. 表结构说明

### 3.1 用户表 (p_user)

适用于单用户博客系统，只保留一个用户记录。

| 字段名 | 类型 | 允许空 | 默认值 | 说明 |
|--------|------|--------|--------|------|
| id | int | 否 | - | 主键 |
| account | varchar(64) | 否 | - | 账户 |
| userName | varchar(64) | 否 | - | 账户名 |
| password | varchar(128) | 否 | - | 密码（加密存储） |
| salt | varchar(64) | 否 | - | 盐值 |
| status | int | 否 | 0 | 状态:0正常,1禁用 |
| avatar | varchar(512) | 是 | NULL | 头像 |
| email | varchar(128) | 是 | NULL | 邮箱 |
| website | varchar(256) | 是 | NULL | 网站地址 |
| createTime | datetime | 否 | - | 创建时间 |
| updateTime | datetime | 是 | NULL | 更新时间 |

### 3.2 文章表 (p_article)

存储博客文章内容及相关信息。

| 字段名 | 类型 | 允许空 | 默认值 | 说明 |
|--------|------|--------|--------|------|
| id | int | 否 | - | 主键 |
| title | varchar(256) | 是 | NULL | 文章标题 |
| content | longtext | 否 | - | 文章内容 |
| parseContent | longtext | 否 | - | 解析后的文章内容 |
| contentModel | varchar(64) | 是 | NULL | 文章内容类型:html/markdown |
| type | varchar(32) | 是 | NULL | 文章类型:article文章,page页面 |
| summary | varchar(1024) | 是 | NULL | 文章摘要 |
| metaKeywords | varchar(512) | 是 | NULL | SEO关键字 |
| metaDescription | varchar(512) | 是 | NULL | SEO描述 |
| thumbnail | varchar(256) | 是 | NULL | 缩略图 |
| slug | varchar(128) | 是 | NULL | slug |
| isTop | int | 是 | 0 | 是否置顶0:否,1:是 |
| status | int | 是 | 0 | 状态0:已发布,1:草稿 |
| viewCount | int | 是 | 0 | 访问量 |
| greatCount | int | 是 | 0 | 点赞数 |
| isComment | int | 是 | 1 | 是否允许评论0:否,1是 |
| flag | varchar(256) | 是 | NULL | 标识 |
| template | varchar(256) | 是 | NULL | 模板 |
| createTime | datetime | 否 | - | 创建时间 |
| updateTime | datetime | 是 | NULL | 更新时间 |
| createUserId | int | 否 | - | 添加人 |
| updateUserId | int | 是 | NULL | 更新人 |
| visibility | int | 否 | 0 | 是否可见, 0是, 1否 |

### 3.3 分类表 (p_category)

文章分类信息。

| 字段名 | 类型 | 允许空 | 默认值 | 说明 |
|--------|------|--------|--------|------|
| id | int | 否 | - | 主键 |
| name | varchar(256) | 否 | - | 分类名 |
| pid | int | 否 | -1 | 父级id |
| desc | varchar(512) | 是 | NULL | 描述 |
| metaKeywords | varchar(256) | 是 | NULL | SEO关键字 |
| thumbnail | varchar(256) | 是 | NULL | 封面图 |
| slug | varchar(128) | 是 | NULL | slug |
| metaDescription | varchar(256) | 是 | NULL | SEO描述内容 |
| status | int | 否 | 0 | 状态0:正常,1禁用 |
| createTime | datetime | 否 | - | 创建时间 |
| updateTime | datetime | 是 | NULL | 更新时间 |
| createUserId | int | 是 | NULL | 添加人 |
| updateUserId | int | 是 | NULL | 更新人 |

### 3.4 标签表 (p_tag)

文章标签信息。

| 字段名 | 类型 | 允许空 | 默认值 | 说明 |
|--------|------|--------|--------|------|
| id | int | 否 | - | 主键 |
| name | varchar(256) | 否 | - | 标签名 |
| color | varchar(128) | 是 | NULL | 颜色 |
| thumbnail | varchar(256) | 是 | NULL | 缩略图 |
| slug | varchar(128) | 是 | NULL | slug |
| createUserId | int | 是 | NULL | 添加人 |
| updateUserId | int | 是 | NULL | 更新人 |
| createTime | datetime | 否 | - | 创建时间 |
| updateTime | datetime | 是 | NULL | 修改时间 |

### 3.5 文章分类关联表 (p_article_category)

文章与分类的多对多关联关系。

| 字段名 | 类型 | 允许空 | 默认值 | 说明 |
|--------|------|--------|--------|------|
| id | int | 否 | - | 主键 |
| articleId | int | 否 | - | 文章id |
| categoryId | int | 否 | - | 分类id |
| createTime | datetime | 否 | - | 创建时间 |
| updateTime | datetime | 是 | NULL | 更新时间 |
| createUserId | int | 是 | NULL | 添加人 |
| updateUserId | int | 是 | NULL | 更新人 |

### 3.6 文章标签关联表 (p_article_tag)

文章与标签的多对多关联关系。

| 字段名 | 类型 | 允许空 | 默认值 | 说明 |
|--------|------|--------|--------|------|
| id | int | 否 | - | 主键 |
| articleId | int | 否 | - | 文章id |
| tagId | int | 否 | - | 标签id |
| createTime | datetime | 是 | NULL | 创建时间 |
| updateTime | datetime | 是 | NULL | 更新时间 |
| createUserId | int | 是 | NULL | 添加人 |
| updateUserId | int | 是 | NULL | 更新人 |

### 3.7 评论表 (p_comment)

文章评论信息，无需登录但需要审核。

| 字段名 | 类型 | 允许空 | 默认值 | 说明 |
|--------|------|--------|--------|------|
| id | int | 否 | - | 主键 |
| articleId | int | 否 | - | 文章id |
| pid | int | 是 | -1 | 父级id |
| topPid | int | 是 | -1 | 顶层父级id |
| content | varchar(2048) | 是 | NULL | 评论内容 |
| status | int | 是 | 0 | 状态:0正常,1:待审核 |
| avatar | varchar(256) | 是 | NULL | 头像 |
| website | varchar(256) | 是 | NULL | 网站地址 |
| email | varchar(256) | 是 | NULL | 邮箱 |
| userName | varchar(256) | 是 | NULL | 评论人 |
| ip | varchar(256) | 是 | NULL | ip |
| device | varchar(256) | 是 | NULL | 设备类型 |
| createTime | datetime | 否 | - | 创建时间 |
| updateTime | datetime | 是 | NULL | 更新时间 |

### 3.8 附件表 (p_attach)

文件上传信息。

| 字段名 | 类型 | 允许空 | 默认值 | 说明 |
|--------|------|--------|--------|------|
| id | int | 否 | - | 主键 |
| name | varchar(256) | 否 | - | 附件名 |
| remark | varchar(512) | 是 | NULL | 附件描述 |
| path | varchar(512) | 否 | - | 附件路径 |
| flag | varchar(256) | 是 | NULL | 标识 |
| mineType | varchar(128) | 是 | NULL | 文件类型mineType |
| type | varchar(32) | 是 | NULL | 文件类型 |
| createTime | datetime | 否 | - | 创建时间 |
| updateTime | datetime | 是 | NULL | 更新时间 |
| url | varchar(512) | 否 | - | 访问路径 |
| attachGroup | varchar(256) | 否 | default | 附件分组 |
| createUserId | int | 是 | NULL | 添加人 |
| updateUserId | int | 是 | NULL | 更新人 |

### 3.9 友情链接表 (p_link)

友情链接信息。

| 字段名 | 类型 | 允许空 | 默认值 | 说明 |
|--------|------|--------|--------|------|
| id | int | 否 | - | 主键 |
| name | varchar(256) | 否 | - | 网站名 |
| logo | varchar(256) | 是 | NULL | 网站logo |
| desc | varchar(512) | 是 | NULL | 网站描述 |
| address | varchar(256) | 否 | - | 网站地址 |
| createTime | datetime | 否 | - | 创建时间 |
| updateTime | datetime | 是 | NULL | 更新时间 |
| createUserId | int | 是 | NULL | 添加人 |
| updateUserId | int | 是 | NULL | 更新人 |

## 4. 优化说明

1. **安全性增强**：
   - 增加了密码字段长度至128字符，支持更安全的加密算法
   - 增加了盐值字段长度至64字符

2. **索引优化**：
   - 为所有表的创建时间字段添加了索引
   - 为关联表的外键字段添加了索引
   - 为常用查询字段添加了索引

3. **字段优化**：
   - 增加了contentModel字段长度至64字符
   - 简化了评论表结构，移除了不必要的userId字段
   - 移除了附件表中不必要的configId字段

4. **简化设计**：
   - 针对单用户博客系统简化了用户表结构
   - 移除了不必要的字段，保持表结构简洁