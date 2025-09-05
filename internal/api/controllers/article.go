package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"matuto-blog/internal/database"
	"matuto-blog/internal/models"
	"matuto-blog/pkg/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ArticleController 文章控制器
type ArticleController struct{}

// ArticleRequest 文章请求结构
type ArticleRequest struct {
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug"`
	Summary    string `json:"summary"`
	Content    string `json:"content" binding:"required"`
	Thumbnail  string `json:"thumbnail"`
	CategoryID uint   `json:"category_id" binding:"required"`
	TagIDs     []uint `json:"tag_ids"`
	IsTop      int8   `json:"is_top"`
	IsComment  int8   `json:"is_comment"`
	Status     int8   `json:"status"`
}

type ArticlePageRequest struct {
	common.PageRequest
	CategoryID uint   `json:"categoryId" form:"categoryId"`
	Title      string `json:"title" form:"title"`
	Status     *int8  `json:"status" form:"status"`
}

// ArticlePage 文章分页
func (a *ArticleController) ArticlePage(c *gin.Context) {
	var pageParam ArticlePageRequest
	if err := c.ShouldBindQuery(&pageParam); err != nil {
		common.BadRequest(c, err.Error())
		return
	}

	if err := pageParam.Validate(); err != nil {
		common.BadRequest(c, err.Error())
		return
	}

	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{})

	if pageParam.Status != nil {
		query = query.Where("status = ?", pageParam.Status)
	}

	if pageParam.CategoryID > 0 {
		query = query.Where("category_id = ?", pageParam.CategoryID)
	}

	if pageParam.Title != "" {
		query = query.Where("title LIKE ?", "%"+pageParam.Title+"%")
	}

	query.Count(&total)

	offset := (pageParam.Page - 1) * pageParam.PageSize
	query.Order("is_top DESC, created_at DESC").
		Limit(pageParam.PageSize).
		Offset(offset).
		Find(&articles)

	common.SuccessPage(c, articles, total, pageParam.Page, pageParam.PageSize)
}

// Index 文章列表页面
func (a *ArticleController) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 10
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	tagID, _ := strconv.Atoi(c.Query("tag_id"))
	keyword := strings.TrimSpace(c.Query("keyword"))

	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{}).
		Preload("Category").
		Preload("Tags").
		Where("status = ?", 1)

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if tagID > 0 {
		query = query.Joins("JOIN article_tags ON articles.id = article_tags.article_id").
			Where("article_tags.tag_id = ?", tagID)
	}

	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("is_top DESC, created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&articles)

	// 获取分类列表
	var categories []models.Category
	database.DB.Find(&categories)

	// 获取标签列表
	var tags []models.Tag
	database.DB.Find(&tags)

	c.HTML(http.StatusOK, "default/index.html", gin.H{
		"articles":   articles,
		"categories": categories,
		"tags":       tags,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (int(total) + pageSize - 1) / pageSize,
		},
		"current_category": categoryID,
		"current_tag":      tagID,
		"keyword":          keyword,
	})
}

// Show 文章详情页面
func (a *ArticleController) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/error.html", gin.H{
			"message": "文章不存在",
		})
		return
	}

	var article models.Article
	if err := database.DB.Preload("Category").
		Preload("Tags").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ? AND parent_id = ?", 1, 0).Order("created_at DESC")
		}).
		Preload("Comments.Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", 1).Order("created_at ASC")
		}).
		Where("id = ? AND status = ?", id, 1).
		First(&article).Error; err != nil {
		c.HTML(http.StatusNotFound, "error/error.html", gin.H{
			"message": "文章不存在",
		})
		return
	}

	// 增加访问量
	database.DB.Model(&article).Update("view_count", gorm.Expr("view_count + ?", 1))
	article.ViewCount++

	c.HTML(http.StatusOK, "article.html", gin.H{
		"article": article,
		"title":   article.Title,
	})
}

// AdminIndex 管理后台文章列表
func (a *ArticleController) AdminIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 15
	status := c.Query("status")
	keyword := strings.TrimSpace(c.Query("keyword"))

	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{}).Preload("Category")

	if status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}

	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&articles)

	c.HTML(http.StatusOK, "admin/articles/index.html", gin.H{
		"articles": articles,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (int(total) + pageSize - 1) / pageSize,
		},
		"current_status": status,
		"keyword":        keyword,
		"title":          "文章管理",
	})
}

// AdminCreate 创建文章页面
func (a *ArticleController) AdminCreate(c *gin.Context) {
	// 获取分类列表
	var categories []models.Category
	database.DB.Where("status = ?", 1).Order("sort ASC").Find(&categories)

	// 获取标签列表
	var tags []models.Tag
	database.DB.Where("status = ?", 1).Order("sort ASC").Find(&tags)

	c.HTML(http.StatusOK, "admin/articles/create.html", gin.H{
		"categories": categories,
		"tags":       tags,
		"title":      "创建文章",
	})
}

// AdminStore 保存文章
func (a *ArticleController) AdminStore(c *gin.Context) {
	var req ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 生成slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Title)
	}

	// 创建文章
	article := models.Article{
		Title:     req.Title,
		Slug:      req.Slug,
		Summary:   req.Summary,
		Content:   req.Content,
		Thumbnail: req.Thumbnail,
		IsTop:     req.IsTop,
		IsComment: req.IsComment,
		Status:    req.Status,
	}

	if req.Status == 1 {
		article.CreatedAt = time.Now()
	}

	if err := database.DB.Create(&article).Error; err != nil {
		common.ServerError(c, "创建文章失败: "+err.Error())
		return
	}

	// 关联标签
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		database.DB.Where("id IN ?", req.TagIDs).Find(&tags)
		database.DB.Model(&article).Association("Tags").Replace(tags)
	}

	common.Success(c, gin.H{
		"id": article.Id,
	})
}

// AdminEdit 编辑文章页面
func (a *ArticleController) AdminEdit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.HTML(http.StatusNotFound, "admin/error/error.html", gin.H{
			"message": "文章不存在",
		})
		return
	}

	var article models.Article
	if err := database.DB.Preload("Tags").First(&article, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "admin/error/error.html", gin.H{
			"message": "文章不存在",
		})
		return
	}

	// 获取分类列表
	var categories []models.Category
	database.DB.Where("status = ?", 1).Order("sort ASC").Find(&categories)

	// 获取标签列表
	var tags []models.Tag
	database.DB.Where("status = ?", 1).Order("sort ASC").Find(&tags)

	// 获取文章已选标签ID
	var selectedTagIDs []int
	for _, tag := range article.Tags {
		selectedTagIDs = append(selectedTagIDs, tag.Id)
	}

	c.HTML(http.StatusOK, "admin/articles/edit.html", gin.H{
		"article":          article,
		"categories":       categories,
		"tags":             tags,
		"selected_tag_ids": selectedTagIDs,
		"title":            "编辑文章",
	})
}

// AdminUpdate 更新文章
func (a *ArticleController) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.BadRequest(c, "无效的文章ID")
		return
	}

	var req ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var article models.Article
	if err := database.DB.First(&article, id).Error; err != nil {
		common.NotFound(c, "文章不存在")
		return
	}

	// 生成slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Title)
	}

	// 更新文章
	article.Title = req.Title
	article.Slug = req.Slug
	article.Summary = req.Summary
	article.Content = req.Content
	article.Thumbnail = req.Thumbnail
	article.IsTop = req.IsTop
	article.IsComment = req.IsComment

	// 如果状态从草稿改为已发布，设置发布时间
	if article.Status == 0 && req.Status == 1 {
		article.CreatedAt = time.Now()
	}
	article.Status = req.Status

	if err := database.DB.Save(&article).Error; err != nil {
		common.ServerError(c, "更新文章失败: "+err.Error())
		return
	}

	// 更新标签关联
	var tags []models.Tag
	if len(req.TagIDs) > 0 {
		database.DB.Where("id IN ?", req.TagIDs).Find(&tags)
	}
	database.DB.Model(&article).Association("Tags").Replace(tags)

	common.SuccessWithMessage(c, "文章更新成功", nil)
}

// AdminDestroy 删除文章
func (a *ArticleController) AdminDestroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.BadRequest(c, "无效的文章ID")
		return
	}

	if err := database.DB.Delete(&models.Article{}, id).Error; err != nil {
		common.ServerError(c, "删除文章失败: "+err.Error())
		return
	}

	common.SuccessWithMessage(c, "文章删除成功", nil)
}

// generateSlug 生成文章slug
func generateSlug(title string) string {
	// 这里简单处理，实际项目中可能需要更复杂的slug生成逻辑
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	return strconv.FormatInt(time.Now().UnixNano(), 10) // 使用时间戳确保唯一性
}
