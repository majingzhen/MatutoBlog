package controllers

import (
	"matuto-blog/internal/database"
	"matuto-blog/internal/models"
	"matuto-blog/pkg/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
type CategoryController struct{}

// CategoryRequest 分类请求结构
type CategoryRequest struct {
	Name            string `json:"name" binding:"required"`
	Pid             int    `json:"pId"`
	Desc            string `json:"desc"`
	Thumbnail       string `json:"thumbnail"`
	Slug            string `json:"slug"`
	MetaKeywords    string `json:"metaKeywords"`
	MetaDescription string `json:"metaDescription"`
	Status          int    `json:"status"`
}

// CategoryPageRequest 分类分页请求
type CategoryPageRequest struct {
	common.PageRequest
	Name   string `json:"name" form:"name"`
	Status *int   `json:"status" form:"status"`
}

// CategoryPage 分类分页
func (c *CategoryController) CategoryPage(ctx *gin.Context) {
	var req CategoryPageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.ServerError(ctx, "参数错误: "+err.Error())
		return
	}

	var categories []models.Category
	var total int64

	query := database.DB.Model(&models.Category{})

	// 名称搜索
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	query.Order("created_at DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&categories)
	common.SuccessPage(ctx, categories, total, req.Page, req.PageSize)
}

// CategoryEnableList 分类启用列表
func (c *CategoryController) CategoryEnableList(ctx *gin.Context) {
	var categories []models.Category
	database.DB.Where("status = ?", models.StatusActive).Order("created_at DESC").Find(&categories)
	common.Success(ctx, categories)
}

// DeleteCategory 删除分类
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		common.ServerError(ctx, "无效的分类ID")
		return
	}

	// 检查是否有子分类
	var childCount int64
	database.DB.Model(&models.Category{}).Where("p_id = ?", id).Count(&childCount)
	if childCount > 0 {
		common.ServerError(ctx, "该分类下有子分类，无法删除")
		return
	}

	// 检查是否有文章
	var articleCount int64
	database.DB.Model(&models.Article{}).Where("category_id = ?", id).Count(&articleCount)
	if articleCount > 0 {
		common.ServerError(ctx, "该分类下有文章，无法删除")
		return
	}

	if err := database.DB.Delete(&models.Category{}, id).Error; err != nil {
		common.ServerError(ctx, "删除分类失败: "+err.Error())
		return
	}
	common.SuccessWithMessage(ctx, "分类删除成功", nil)
}

// CreateCategory 创建分类
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var req CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.ServerError(ctx, "参数错误: "+err.Error())
		return
	}

	// 生成slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	category := models.Category{
		Name:            req.Name,
		Pid:             req.Pid,
		Desc:            req.Desc,
		Thumbnail:       req.Thumbnail,
		Slug:            req.Slug,
		MetaKeywords:    req.MetaKeywords,
		MetaDescription: req.MetaDescription,
		Status:          req.Status,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		common.ServerError(ctx, "创建分类失败: "+err.Error())
		return
	}
	common.SuccessWithMessage(ctx, "分类创建成功", nil)
}

// UpdateCategory 更新分类
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		common.ServerError(ctx, "无效的分类ID")
		return
	}

	var req CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.ServerError(ctx, "参数错误: "+err.Error())
		return
	}

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		common.ServerError(ctx, "分类不存在")
		return
	}

	// 检查父分类不能是自己或子分类
	if req.Pid == int(id) {
		common.ServerError(ctx, "父分类不能是自己")
		return
	}

	// 生成slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	// 更新字段
	category.Name = req.Name
	category.Pid = req.Pid
	category.Desc = req.Desc
	category.Thumbnail = req.Thumbnail
	category.Slug = req.Slug
	category.MetaKeywords = req.MetaKeywords
	category.MetaDescription = req.MetaDescription
	category.Status = req.Status

	if err := database.DB.Save(&category).Error; err != nil {
		common.ServerError(ctx, "更新分类失败: "+err.Error())
		return
	}
	common.SuccessWithMessage(ctx, "分类更新成功", nil)
}
