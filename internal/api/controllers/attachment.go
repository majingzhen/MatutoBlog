package controllers

import (
	"fmt"
	"io"
	"matuto-blog/internal/database"
	"matuto-blog/internal/models"
	"matuto-blog/pkg/common"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// AttachmentController 附件控制器
type AttachmentController struct{}

type AttachPageRequest struct {
	common.PageRequest
	Title string `json:"title" form:"title"`
}

// AttachPage 附件列表
func (r *AttachmentController) AttachPage(ctx *gin.Context) {
	var req AttachPageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	var attachments []models.Attach
	var total int64

	query := database.DB.Model(&models.Attach{})

	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	query.Order("created_at DESC").
		Limit(req.PageSize).
		Offset(offset).
		Find(&attachments)

	common.SuccessPage(ctx, attachments, total, req.Page, req.PageSize)
}

// Upload 文件上传
func (a *AttachmentController) Upload(ctx *gin.Context) {
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		common.ServerError(ctx, "请选择要上传的文件")
		return
	}

	// 检查文件大小
	maxSize := viper.GetInt64("storage.local.max_size")
	if maxSize == 0 {
		maxSize = 10 * 1024 * 1024 // 默认10MB
	}
	if file.Size > maxSize {
		common.ServerError(ctx, "文件大小超过限制")
		return
	}

	// 检查文件类型
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".txt", ".zip", ".rar"}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		common.ServerError(ctx, "不支持的文件类型")
		return
	}

	// 创建上传目录
	uploadPath := viper.GetString("storage.local.base_path")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	// 按年月分目录
	now := time.Now()
	datePath := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
	fullPath := filepath.Join(uploadPath, datePath)

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		common.ServerError(ctx, "创建上传目录失败")
		return
	}

	// 生成新文件名
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(fullPath, filename)

	// 保存文件
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		common.ServerError(ctx, "保存文件失败: "+err.Error())
		return
	}

	// 获取文件MIME类型
	src, err := file.Open()
	if err != nil {
		os.Remove(filePath) // 删除已保存的文件
		common.ServerError(ctx, "读取文件失败: "+err.Error())
		return
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		os.Remove(filePath)
		common.ServerError(ctx, "读取文件失败: "+err.Error())
		return
	}

	// 保存到数据库
	attachment := models.Attach{
		Name: file.Filename,
		Path: strings.ReplaceAll(filepath.Join(datePath, filename), "\\", "/"),
	}

	if err := database.DB.Create(&attachment).Error; err != nil {
		os.Remove(filePath) // 删除已保存的文件
		common.ServerError(ctx, "保存文件记录失败: "+err.Error())
		return
	}

	// 生成访问URL
	baseURL := viper.GetString("storage.local.base_url")
	if baseURL == "" {
		baseURL = "http://localhost:8080/uploads/"
	}
	fileURL := baseURL + attachment.Path

	common.SuccessWithMessage(ctx, "文件上传成功", gin.H{
		"id":   attachment.Id,
		"name": attachment.Name,
		"url":  fileURL,
	})
}

// AdminIndex 管理后台附件列表
func (a *AttachmentController) AdminIndex(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize := 20
	keyword := strings.TrimSpace(ctx.Query("keyword"))

	var attachments []models.Attach
	var total int64

	query := database.DB.Model(&models.Attach{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&attachments)

	// 为每个附件生成访问URL
	baseURL := viper.GetString("storage.local.base_url")
	if baseURL == "" {
		baseURL = "http://localhost:8080/uploads/"
	}

	for i := range attachments {
		attachments[i].Path = baseURL + attachments[i].Path
	}

	ctx.HTML(http.StatusOK, "admin/attachments/index.html", gin.H{
		"attachments": attachments,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"pages":     (int(total) + pageSize - 1) / pageSize,
		},
		"keyword": keyword,
		"title":   "附件管理",
	})
}

// AdminDestroy 删除附件
func (a *AttachmentController) AdminDestroy(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的附件ID",
		})
		return
	}

	var attachment models.Attach
	if err := database.DB.First(&attachment, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "附件不存在",
		})
		return
	}

	// 删除物理文件
	uploadPath := viper.GetString("storage.local.base_path")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	filePath := filepath.Join(uploadPath, strings.ReplaceAll(attachment.Path, "/", string(os.PathSeparator)))
	if err := os.Remove(filePath); err != nil {
		// 文件不存在也继续删除数据库记录
		fmt.Printf("删除文件失败: %v\n", err)
	}

	// 删除数据库记录
	if err := database.DB.Delete(&attachment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除附件记录失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "附件删除成功",
	})
}

// AdminBatchDelete 批量删除附件
func (a *AttachmentController) AdminBatchDelete(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	if len(req.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请选择要删除的附件",
		})
		return
	}

	// 获取要删除的附件信息
	var attachments []models.Attach
	database.DB.Where("id IN ?", req.IDs).Find(&attachments)

	uploadPath := viper.GetString("storage.local.base_path")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	// 删除物理文件
	for _, attachment := range attachments {
		filePath := filepath.Join(uploadPath, strings.ReplaceAll(attachment.Path, "/", string(os.PathSeparator)))
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("删除文件失败: %v\n", err)
		}
	}

	// 批量删除数据库记录
	if err := database.DB.Where("id IN ?", req.IDs).Delete(&models.Attach{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "批量删除失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("已删除 %d 个附件", len(attachments)),
	})
}

// GetUploadToken 获取上传令牌（如果需要的话）
func (a *AttachmentController) GetUploadToken(ctx *gin.Context) {
	// 这里可以实现上传令牌生成逻辑
	// 暂时返回简单的成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取上传令牌成功",
		"data": gin.H{
			"token":      "simple-upload-token",
			"expires_at": time.Now().Add(time.Hour).Unix(),
		},
	})
}
