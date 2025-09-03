package controllers

import (
	"fmt"
	"matuto-blog/internal/database"
	"matuto-blog/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CommentController 评论控制器
type CommentController struct{}

// CommentRequest 评论请求结构
type CommentRequest struct {
	ArticleID int    `json:"articleId" form:"articleId" binding:"required"`
	Pid       int    `json:"pId" form:"pId"`
	UserName  string `json:"username" form:"userName" binding:"required"`
	Email     string `json:"email" form:"email"`
	Website   string `json:"website" form:"website"`
	Content   string `json:"content" form:"content" binding:"required"`
}

// Submit 提交评论
func (c *CommentController) Submit(ctx *gin.Context) {
	var req CommentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		// 如果是AJAX请求，返回JSON
		if ctx.GetHeader("X-Requested-With") == "XMLHttpRequest" || ctx.GetHeader("Content-Type") == "application/json" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "参数错误: " + err.Error(),
			})
			return
		}

		// 否则重定向回文章页面
		ctx.Redirect(http.StatusFound, "/article/"+strconv.Itoa(int(req.ArticleID)))
		return
	}

	// 验证文章是否存在且允许评论
	var article models.Article
	if err := database.DB.Where("id = ? AND status = ? AND is_comment = ?", req.ArticleID, 1, true).First(&article).Error; err != nil {
		if ctx.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "文章不存在或不允许评论",
			})
			return
		}
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	// 获取客户端IP和User-Agent
	clientIP := ctx.ClientIP()
	device := ctx.GetHeader("User-Agent")

	// 创建评论
	comment := models.Comment{
		ArticleId: req.ArticleID,
		Pid:       req.Pid,
		Username:  req.UserName,
		Email:     req.Email,
		Website:   req.Website,
		Content:   req.Content,
		Ip:        clientIP,
		Device:    device,
		Status:    0, // 待审核
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		if ctx.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "评论提交失败: " + err.Error(),
			})
			return
		}
		ctx.Redirect(http.StatusFound, "/article/"+strconv.Itoa(int(req.ArticleID)))
		return
	}

	if ctx.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "评论提交成功，请等待审核",
		})
	} else {
		ctx.Redirect(http.StatusFound, "/article/"+strconv.Itoa(int(req.ArticleID))+"?comment=success")
	}
}

// AdminIndex 管理后台评论列表
func (c *CommentController) AdminIndex(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize := 20
	status := ctx.Query("status")
	keyword := strings.TrimSpace(ctx.Query("keyword"))

	var comments []models.Comment
	var total int64
	query := database.DB.Model(&models.Comment{}).Preload("Article", "id, title")

	// 状态筛选
	if status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("user_name LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&comments)

	// 获取统计数据
	var pendingCount, approvedCount, rejectedCount int64
	database.DB.Model(&models.Comment{}).Where("status = ?", 0).Count(&pendingCount)
	database.DB.Model(&models.Comment{}).Where("status = ?", 1).Count(&approvedCount)
	database.DB.Model(&models.Comment{}).Where("status = ?", 2).Count(&rejectedCount)

	ctx.HTML(http.StatusOK, "admin/comments/index.html", gin.H{
		"comments": comments,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (int(total) + pageSize - 1) / pageSize,
		},
		"current_status": status,
		"keyword":        keyword,
		"stats": gin.H{
			"pending":  pendingCount,
			"approved": approvedCount,
			"rejected": rejectedCount,
		},
		"title": "评论管理",
	})
}

// AdminReview 审核评论
func (c *CommentController) AdminReview(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的评论ID",
		})
		return
	}

	status, err := strconv.Atoi(ctx.PostForm("status"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的状态值",
		})
		return
	}

	// 验证状态值
	if status < 0 || status > 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "状态值必须是 0(待审核), 1(已通过), 2(已拒绝)",
		})
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "评论不存在",
		})
		return
	}

	// 更新评论状态
	comment.Status = status
	if err := database.DB.Save(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新评论状态失败: " + err.Error(),
		})
		return
	}

	// 如果评论通过审核，更新文章评论数量
	if status == 1 {
		var total int64
		database.DB.Model(&models.Article{}).Where("id = ?", comment.ArticleId).Update("comment_count", database.DB.Model(&models.Comment{}).Where("article_id = ? AND status = ?", comment.ArticleId, 1).Count(&total))
	}

	var statusText string
	switch status {
	case 0:
		statusText = "待审核"
	case 1:
		statusText = "已通过"
	case 2:
		statusText = "已拒绝"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "评论状态已更新为: " + statusText,
	})
}

// AdminDestroy 删除评论
func (c *CommentController) AdminDestroy(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的评论ID",
		})
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "评论不存在",
		})
		return
	}

	articleID := comment.ArticleId

	// 删除评论及其回复
	if err := database.DB.Where("id = ? OR parent_id = ?", id, id).Delete(&models.Comment{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除评论失败: " + err.Error(),
		})
		return
	}

	// 更新文章评论数量
	var total int64
	database.DB.Model(&models.Comment{}).Where("article_id = ? AND status = ?", articleID, 1).Count(&total)
	database.DB.Model(&models.Article{}).Where("id = ?", articleID).Update("comment_count", total)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "评论删除成功",
	})
}

// AdminBatchReview 批量审核评论
func (c *CommentController) AdminBatchReview(ctx *gin.Context) {
	var req struct {
		IDs    []uint `json:"ids" binding:"required"`
		Status int    `json:"status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 验证状态值
	if req.Status < 0 || req.Status > 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "状态值必须是 0(待审核), 1(已通过), 2(已拒绝)",
		})
		return
	}

	if len(req.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请选择要操作的评论",
		})
		return
	}

	// 批量更新状态
	if err := database.DB.Model(&models.Comment{}).Where("id IN ?", req.IDs).Update("status", req.Status).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "批量更新失败: " + err.Error(),
		})
		return
	}

	// 如果是通过审核，需要更新相关文章的评论数量
	if req.Status == 1 {
		var comments []models.Comment
		database.DB.Where("id IN ?", req.IDs).Find(&comments)

		// 按文章ID分组更新评论数量
		articleIDs := make(map[int]bool)
		for _, comment := range comments {
			articleIDs[comment.ArticleId] = true
		}

		for articleID := range articleIDs {
			var total int64
			database.DB.Model(&models.Comment{}).Where("article_id = ? AND status = ?", articleID, 1).Count(&total)
			database.DB.Model(&models.Article{}).Where("id = ?", articleID).Update("comment_count", total)
		}
	}

	var statusText string
	switch req.Status {
	case 0:
		statusText = "待审核"
	case 1:
		statusText = "已通过"
	case 2:
		statusText = "已拒绝"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("已将 %d 条评论状态更新为: %s", len(req.IDs), statusText),
	})
}
