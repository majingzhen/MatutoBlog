package controllers

import (
	"matuto-blog/internal/database"
	"matuto-blog/internal/models"
	"matuto-blog/pkg/common"
	"matuto-blog/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
type CategoryController struct{}

// CategoryRequest 分类请求结构
type CategoryRequest struct {
	Id              int    `json:"id"`
	PId             int    `json:"pId"`
	Name            string `json:"name" binding:"required"`
	Thumbnail       string `json:"thumbnail"`
	Slug            string `json:"slug"`
	Desc            string `json:"desc"`
	MetaKeywords    string `json:"metaKeywords"`
	MetaDescription string `json:"metaDescription"`
}

// CategoryPageRequest 分类分页请求
type CategoryPageRequest struct {
	common.PageRequest
	Name string `json:"name" form:"name"`
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

	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	query.Order("created_at DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&categories)
	common.SuccessPage(ctx, categories, total, req.Page, req.PageSize)
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

// AddCategory 添加分类
func (c *CategoryController) AddCategory(ctx *gin.Context) {
	var req CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.ServerError(ctx, "参数错误: "+err.Error())
		return
	}
	to, err := utils.ConvertTo[models.Category](req)
	if err != nil {
		common.ServerError(ctx, "参数错误: "+err.Error())
		return
	}
	// 生成slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	if err := database.DB.Create(&to).Error; err != nil {
		common.ServerError(ctx, "添加分类失败: "+err.Error())
		return
	}
}

// UpdateCategory 更新分类
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		common.ServerError(ctx, "无效的分类ID")
		return
	}

	var req CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.ServerError(ctx, "参数错误: "+err.Error())
		return
	}
	to, err := utils.ConvertTo[models.Category](req)
	if err != nil {
		common.ServerError(ctx, "参数错误: "+err.Error())
		return
	}
	if err := database.DB.First(&models.Category{}, id).Error; err != nil {
		common.ServerError(ctx, "分类不存在")
		return
	}

	// 检查父分类不能是自己或子分类
	if req.PId == int(id) {
		common.ServerError(ctx, "父分类不能是自己或子分类")
		return
	}

	// 生成slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}
	if err := database.DB.Save(&to).Error; err != nil {
		common.ServerError(ctx, "更新分类失败: "+err.Error())
		return
	}
	common.SuccessWithMessage(ctx, "分类更新成功", nil)
}
