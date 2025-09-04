package controllers

import (
	"matuto-blog/internal/database"
	"matuto-blog/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
type CategoryController struct{}

// CategoryRequest 分类请求结构
type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    uint   `json:"parent_id"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
}

// AdminIndex 管理后台分类列表
func (c *CategoryController) AdminIndex(ctx *gin.Context) {
	var categories []models.Category

	// 获取所有分类，包括父子关系
	database.DB.Preload("Children").Where("parent_id = ?", 0).Order("sort ASC").Find(&categories)

	ctx.HTML(http.StatusOK, "admin/categories/index.html", gin.H{
		"categories": categories,
		"title":      "分类管理",
	})
}

// AdminCreate 创建分类页面
func (c *CategoryController) AdminCreate(ctx *gin.Context) {
	// 获取父分类列表
	var parentCategories []models.Category
	database.DB.Where("parent_id = ? AND status = ?", 0, 1).Order("sort ASC").Find(&parentCategories)

	ctx.HTML(http.StatusOK, "admin/categories/create.html", gin.H{
		"parent_categories": parentCategories,
		"title":             "创建分类",
	})
}

// AdminStore 保存分类
func (c *CategoryController) AdminStore(ctx *gin.Context) {
	var req CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 生成slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	// 创建分类
	category := models.Category{
		Name:   req.Name,
		Slug:   req.Slug,
		Status: req.Status,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建分类失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类创建成功",
		"data": gin.H{
			"id": category.Id,
		},
	})
}

// AdminEdit 编辑分类页面
func (c *CategoryController) AdminEdit(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "admin/error/error.html", gin.H{
			"message": "分类不存在",
		})
		return
	}

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		ctx.HTML(http.StatusNotFound, "admin/error/error.html", gin.H{
			"message": "分类不存在",
		})
		return
	}

	// 获取父分类列表（排除自己和子分类）
	var parentCategories []models.Category
	database.DB.Where("parent_id = ? AND status = ? AND id != ?", 0, 1, id).Order("sort ASC").Find(&parentCategories)

	ctx.HTML(http.StatusOK, "admin/categories/edit.html", gin.H{
		"category":          category,
		"parent_categories": parentCategories,
		"title":             "编辑分类",
	})
}

// AdminUpdate 更新分类
func (c *CategoryController) AdminUpdate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的分类ID",
		})
		return
	}

	var req CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "分类不存在",
		})
		return
	}

	// 检查父分类不能是自己或子分类
	if req.ParentID == uint(id) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "父分类不能是自己",
		})
		return
	}

	// 生成slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	// 更新分类
	category.Name = req.Name
	category.Slug = req.Slug
	category.Status = req.Status

	if err := database.DB.Save(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新分类失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类更新成功",
	})
}

// AdminDestroy 删除分类
func (c *CategoryController) AdminDestroy(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的分类ID",
		})
		return
	}

	// 检查是否有子分类
	var childCount int64
	database.DB.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该分类下有子分类，无法删除",
		})
		return
	}

	// 检查是否有文章
	var articleCount int64
	database.DB.Model(&models.Article{}).Where("category_id = ?", id).Count(&articleCount)
	if articleCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该分类下有文章，无法删除",
		})
		return
	}

	if err := database.DB.Delete(&models.Category{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除分类失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类删除成功",
	})
}

// TagController 标签控制器
type TagController struct{}

// TagRequest 标签请求结构
type TagRequest struct {
	Name   string `json:"name" binding:"required"`
	Color  string `json:"color"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status"`
}

// AdminIndex 管理后台标签列表
func (t *TagController) AdminIndex(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize := 20
	keyword := strings.TrimSpace(ctx.Query("keyword"))

	var tags []models.Tag
	var total int64

	query := database.DB.Model(&models.Tag{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("sort ASC, created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&tags)

	ctx.HTML(http.StatusOK, "admin/tags/index.html", gin.H{
		"tags": tags,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (int(total) + pageSize - 1) / pageSize,
		},
		"keyword": keyword,
		"title":   "标签管理",
	})
}

// AdminCreate 创建标签页面
func (t *TagController) AdminCreate(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/tags/create.html", gin.H{
		"title": "创建标签",
	})
}

// AdminStore 保存标签
func (t *TagController) AdminStore(ctx *gin.Context) {
	var req TagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认颜色
	if req.Color == "" {
		req.Color = "#007bff"
	}

	// 创建标签
	tag := models.Tag{
		Name:  req.Name,
		Color: req.Color,
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建标签失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "标签创建成功",
		"data": gin.H{
			"id": tag.Id,
		},
	})
}

// AdminEdit 编辑标签页面
func (t *TagController) AdminEdit(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "admin/error/error.html", gin.H{
			"message": "标签不存在",
		})
		return
	}

	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		ctx.HTML(http.StatusNotFound, "admin/error/error.html", gin.H{
			"message": "标签不存在",
		})
		return
	}

	ctx.HTML(http.StatusOK, "admin/tags/edit.html", gin.H{
		"tag":   tag,
		"title": "编辑标签",
	})
}

// AdminUpdate 更新标签
func (t *TagController) AdminUpdate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的标签ID",
		})
		return
	}

	var req TagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "标签不存在",
		})
		return
	}

	// 设置默认颜色
	if req.Color == "" {
		req.Color = "#007bff"
	}

	// 更新标签
	tag.Name = req.Name
	tag.Color = req.Color

	if err := database.DB.Save(&tag).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新标签失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "标签更新成功",
	})
}

// AdminDestroy 删除标签
func (t *TagController) AdminDestroy(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的标签ID",
		})
		return
	}

	// 检查是否有关联的文章
	var count int64
	database.DB.Table("article_tags").Where("tag_id = ?", id).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该标签下有关联文章，无法删除",
		})
		return
	}

	if err := database.DB.Delete(&models.Tag{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除标签失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "标签删除成功",
	})
}
