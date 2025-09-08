package controllers

import (
	"matuto-blog/internal/database"
	"matuto-blog/internal/models"
	"matuto-blog/pkg/common"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TagController 标签控制器
type TagController struct{}

// TagRequest 标签请求结构
type TagRequest struct {
	Name      string `json:"name" binding:"required"`
	Color     string `json:"color"`
	Slug      string `json:"slug"`
	Thumbnail string `json:"thumbnail"`
}

// TagPageRequest 标签分页请求结构
type TagPageRequest struct {
	common.PageRequest
	Name string `json:"name" form:"name"`
}

// TagPage 标签分页
func (t *TagController) TagPage(ctx *gin.Context) {
	var req TagPageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.ServerError(ctx, "参数错误: "+err.Error())
		return
	}
	var tags []models.Tag
	var total int64
	query := database.DB.Model(&models.Tag{})
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	query.Count(&total)
	offset := (req.Page - 1) * req.PageSize
	query.Order("sort ASC, created_at DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&tags)
	common.SuccessPage(ctx, tags, total, req.Page, req.PageSize)
}

// DeleteTag 删除标签
func (t *TagController) DeleteTag(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		common.ServerError(ctx, "参数错误: "+err.Error())
		return
	}
	// 检查是否有关联的文章
	var count int64
	database.DB.Table("article_tags").Where("tag_id = ?", id).Count(&count)
	if count > 0 {
		common.ServerError(ctx, "该标签下有关联文章，无法删除")
		return
	}

	if err := database.DB.Delete(&models.Tag{}, id).Error; err != nil {
		common.ServerError(ctx, "删除标签失败: "+err.Error())
		return
	}

	common.SuccessWithMessage(ctx, "标签删除成功", nil)
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
