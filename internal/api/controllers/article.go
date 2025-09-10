package controllers

import (
	"matuto-blog/pkg/utils"
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
	Id              int      `json:"id"`
	Title           string   `json:"title" binding:"required"`
	Slug            string   `json:"slug"`
	Summary         string   `json:"summary"`
	Content         string   `json:"content" binding:"required"`
	Thumbnail       string   `json:"thumbnail"`
	CategoryIds     []int    `json:"categoryIds"`
	MetaTitle       string   `json:"metaTitle"`
	MetaKeywords    string   `json:"metaKeywords"`
	MetaDescription string   `json:"metaDescription"`
	ContentModel    string   `json:"contentModel"`
	Type            string   `json:"type"`
	TagIDs          []int    `json:"tagIds"`
	AddTags         []string `json:"addTags"`
	IsTop           int8     `json:"isTop"`
	IsComment       int8     `json:"isComment"`
	Status          int8     `json:"status"`
}

// ArticleResponse 文章响应结构
type ArticleResponse struct {
	models.Article
	CategoryIds []int `json:"categoryIds"`
	TagIds      []int `json:"tagIds"`
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

// GetArticle 获取文章详情
func (a *ArticleController) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.BadRequest(c, "无效的文章ID")
		return
	}

	var article models.Article
	if err := database.DB.Where("id = ?", id).First(&article).Error; err != nil {
		common.ServerError(c, "文章不存在")
		return
	}

	// 查询关联分类ID，不使用关联查询
	var categoryIds []int
	database.DB.Model(&models.ArticleCategory{}).
		Where("article_id = ?", id).
		Pluck("category_id", &categoryIds)

	// 查询关联标签ID，不使用关联查询
	var tagIds []int
	database.DB.Model(&models.ArticleTag{}).
		Where("article_id = ?", id).
		Pluck("tag_id", &tagIds)

	response := ArticleResponse{
		Article:     article,
		CategoryIds: categoryIds,
		TagIds:      tagIds,
	}

	common.Success(c, response)
}

// DeleteArticle 删除文章
func (a *ArticleController) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		common.BadRequest(c, "无效的文章ID")
		return
	}

	// 开启事务删除文章及其关联
	tx := database.DB.Begin()

	// 删除文章分类关联
	if err := tx.Where("article_id = ?", id).Delete(&models.ArticleCategory{}).Error; err != nil {
		tx.Rollback()
		common.ServerError(c, "删除文章分类关联失败: "+err.Error())
		return
	}

	// 删除文章标签关联
	if err := tx.Where("article_id = ?", id).Delete(&models.ArticleTag{}).Error; err != nil {
		tx.Rollback()
		common.ServerError(c, "删除文章标签关联失败: "+err.Error())
		return
	}

	// 删除文章
	if err := tx.Delete(&models.Article{}, id).Error; err != nil {
		tx.Rollback()
		common.ServerError(c, "删除文章失败: "+err.Error())
		return
	}

	tx.Commit()
	common.SuccessWithMessage(c, "文章删除成功", nil)
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

// PublishArticle 发布文章
func (a *ArticleController) PublishArticle(c *gin.Context) {
	var req ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	tagIds := req.TagIDs

	// 生成slug
	if req.Slug == "" {
		req.Slug = utils.GenerateSlug(req.Title)
	}
	// 如果addTag不为空，则创建标签
	if len(req.AddTags) > 0 {
		for _, tagName := range req.AddTags {
			tag := models.Tag{
				Name: tagName,
			}
			if err := database.DB.Where("name = ?", tagName).FirstOrCreate(&tag).Error; err != nil {
				common.ServerError(c, "创建标签失败: "+err.Error())
				return
			}
			tagIds = append(tagIds, tag.Id)
		}
	}

	// 创建文章
	article, err := utils.ConvertTo[models.Article](req)
	if err != nil {
		common.ServerError(c, "参数错误: "+err.Error())
		return
	}

	if req.Status == 1 {
		article.CreatedAt = time.Now()
	}

	if err := database.DB.Create(&article).Error; err != nil {
		common.ServerError(c, "创建文章失败: "+err.Error())
		return
	}

	// 关联标签
	if len(tagIds) > 0 {
		var tags []models.ArticleTag
		for _, tagID := range tagIds {
			tags = append(tags, models.ArticleTag{
				ArticleId: article.Id,
				TagId:     tagID,
			})
		}
		if err := database.DB.Create(&tags).Error; err != nil {
			common.ServerError(c, "关联标签失败: "+err.Error())
			return
		}
	}
	// 关联分类
	if len(req.CategoryIds) > 0 {
		var categories []models.ArticleCategory
		for _, categoryID := range req.CategoryIds {
			categories = append(categories, models.ArticleCategory{
				ArticleId:  article.Id,
				CategoryId: categoryID,
			})
		}
		if err := database.DB.Create(&categories).Error; err != nil {
			common.ServerError(c, "关联分类失败: "+err.Error())
			return
		}
	}

	common.Success(c, gin.H{
		"id": article.Id,
	})
}

// UpdateArticle 更新文章
func (a *ArticleController) UpdateArticle(c *gin.Context) {
	var req ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := database.DB.First(&models.Article{}, req.Id).Error; err != nil {
		common.NotFound(c, "文章不存在")
		return
	}

	// 生成slug
	if req.Slug == "" {
		req.Slug = utils.GenerateSlug(req.Title)
	}

	// 更新文章
	article, err := utils.ConvertTo[models.Article](req)
	if err != nil {
		common.ServerError(c, "参数错误: "+err.Error())
		return
	}
	// 如果状态从草稿改为已发布，设置发布时间
	if article.Status == 0 && req.Status == 1 {
		article.CreatedAt = time.Now()
	}
	article.Status = req.Status

	if err := database.DB.Save(&article).Error; err != nil {
		common.ServerError(c, "更新文章失败: "+err.Error())
		return
	}
	tagIds := req.TagIDs
	// 如果addTag不为空，则创建标签
	if len(req.AddTags) > 0 {
		for _, tagName := range req.AddTags {
			tag := models.Tag{
				Name: tagName,
			}
			if err := database.DB.Where("name = ?", tagName).FirstOrCreate(&tag).Error; err != nil {
				common.ServerError(c, "创建标签失败: "+err.Error())
				return
			}
			tagIds = append(tagIds, tag.Id)
		}
	}
	// 删除旧的标签关联
	if err := database.DB.Where("article_id = ?", req.Id).Delete(&models.ArticleTag{}).Error; err != nil {
		common.ServerError(c, "删除旧标签关联失败: "+err.Error())
		return
	}
	// 删除旧的分类关联
	if err := database.DB.Where("article_id = ?", req.Id).Delete(&models.ArticleCategory{}).Error; err != nil {
		common.ServerError(c, "删除旧分类关联失败: "+err.Error())
		return
	}

	// 关联标签
	if len(tagIds) > 0 {
		var tags []models.ArticleTag
		for _, tagID := range tagIds {
			tags = append(tags, models.ArticleTag{
				ArticleId: article.Id,
				TagId:     tagID,
			})
		}
		if err := database.DB.Create(&tags).Error; err != nil {
			common.ServerError(c, "关联标签失败: "+err.Error())
			return
		}
	}
	// 关联分类
	if len(req.CategoryIds) > 0 {
		var categories []models.ArticleCategory
		for _, categoryID := range req.CategoryIds {
			categories = append(categories, models.ArticleCategory{
				ArticleId:  article.Id,
				CategoryId: categoryID,
			})
		}
		if err := database.DB.Create(&categories).Error; err != nil {
			common.ServerError(c, "关联分类失败: "+err.Error())
			return
		}
	}

	common.SuccessWithMessage(c, "文章更新成功", nil)
}
