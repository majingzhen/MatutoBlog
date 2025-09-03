package models

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// Attach 附件模型
type Attach struct {
	BaseModel
	Name   string `json:"name" gorm:"size:256;not null;comment:附件名"`
	Remark string `json:"remark" gorm:"size:512;comment:附件描述"`
	Path   string `json:"path" gorm:"size:512;not null;comment:附件路径"`
	Flag   string `json:"flag" gorm:"size:256;comment:标识"`
	Type   string `json:"type" gorm:"size:32;index;comment:文件类型"`
	URL    string `json:"url" gorm:"size:512;not null;comment:访问路径"`
}

// TableName 指定表名
func (Attach) TableName() string {
	return "m_attach"
}

// AttachType 附件类型常量
const (
	AttachTypeImage = "image" // 图片
	AttachTypeFile  = "file"  // 文件
	AttachTypeVideo = "video" // 视频
	AttachTypeAudio = "audio" // 音频
)

// StorageType 存储类型常量
const (
	StorageTypeLocal = 0 // 本地存储
	StorageTypeOSS   = 1 // OSS存储
	StorageTypeCOS   = 2 // COS存储
)

// 支持的图片格式
var SupportedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// 支持的文件扩展名
var SupportedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// 文件大小限制（字节）
const (
	MaxImageSize   = 10 * 1024 * 1024 // 10MB
	MaxGeneralSize = 50 * 1024 * 1024 // 50MB
)

// GetSizeString 获取文件大小的可读字符串
func (a *Attach) GetSizeString() string {
	// Note: 由于数据库中没有存储文件大小，这里无法实现真实的大小显示
	return "未知大小"
}

// ValidateImageFile 验证图片文件
func ValidateImageFile(filename string, size int64, mimeType string) error {
	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))
	if !SupportedImageExtensions[ext] {
		return fmt.Errorf("不支持的图片格式，仅支持: jpg, jpeg, png, gif, webp")
	}

	// 检查MIME类型
	if mimeType != "" && !SupportedImageTypes[mimeType] {
		return fmt.Errorf("不支持的图片类型: %s", mimeType)
	}

	// 检查文件大小
	if size > MaxImageSize {
		return fmt.Errorf("图片文件大小不能超过 %s", getSizeString(MaxImageSize))
	}

	return nil
}

// ValidateGeneralFile 验证通用文件
func ValidateGeneralFile(filename string, size int64) error {
	// 检查文件大小
	if size > MaxGeneralSize {
		return fmt.Errorf("文件大小不能超过 %s", getSizeString(MaxGeneralSize))
	}

	// 检查文件名
	if len(filename) > 255 {
		return fmt.Errorf("文件名长度不能超过255个字符")
	}

	return nil
}

// getSizeString 获取大小的可读字符串
func getSizeString(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// GenerateAttachPath 生成附件存储路径
func GenerateAttachPath(category, filename string) string {
	// 获取文件扩展名
	ext := filepath.Ext(filename)

	// 生成时间路径
	now := time.Now()
	datePath := now.Format("2006/01/02")

	// 生成唯一文件名
	timestamp := now.Unix()
	uniqueName := fmt.Sprintf("%d_%s%s", timestamp, generateRandomString(8), ext)

	// 组合完整路径: category/2006/01/02/timestamp_random.ext
	return fmt.Sprintf("%s/%s/%s", category, datePath, uniqueName)
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		// 简单的随机数生成方式
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
