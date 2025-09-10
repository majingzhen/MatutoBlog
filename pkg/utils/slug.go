package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// GenerateSlug 生成 slug，支持中文转拼音首字母和英文转小写
// 中文字符会转换为拼音首字母，英文字符转换为小写，其他字符会被移除或替换为连字符
func GenerateSlug(input string) string {
	if input == "" {
		return ""
	}

	var result strings.Builder

	// 遍历每个字符
	for _, char := range input {
		if unicode.Is(unicode.Han, char) {
			// 中文字符：获取拼音首字母
			pinyinResult := pinyin.LazyConvert(string(char), nil)
			if len(pinyinResult) > 0 && len(pinyinResult[0]) > 0 {
				// 取拼音首字母并转小写
				firstLetter := strings.ToLower(string(pinyinResult[0][0]))
				result.WriteString(firstLetter)
			}
		} else if unicode.IsLetter(char) {
			// 英文字母：转小写
			result.WriteString(strings.ToLower(string(char)))
		} else if unicode.IsDigit(char) {
			// 数字：保持不变
			result.WriteString(string(char))
		} else if unicode.IsSpace(char) || char == '-' || char == '_' {
			// 空格、连字符、下划线：转换为连字符
			result.WriteString("-")
		}
		// 其他字符忽略
	}

	slug := result.String()

	// 清理多余的连字符
	slug = cleanupSlug(slug)

	return slug
}

// cleanupSlug 清理 slug 中的多余连字符
func cleanupSlug(slug string) string {
	// 移除开头和结尾的连字符
	slug = strings.Trim(slug, "-")

	// 将多个连续的连字符替换为单个连字符
	re := regexp.MustCompile("-+")
	slug = re.ReplaceAllString(slug, "-")

	return slug
}
