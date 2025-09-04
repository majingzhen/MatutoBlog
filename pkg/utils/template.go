package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// TemplateManager 管理模板相关配置
type TemplateManager struct {
	templateDir  string
	templateDirs []string
}

// NewTemplateManager 创建新的模板管理器
func NewTemplateManager(templateDir string, templateNames []string) *TemplateManager {
	return &TemplateManager{
		templateDir:  templateDir,
		templateDirs: templateNames,
	}
}

// GenTemplateFuncMap 添加自定义模板函数
func GenTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"substr": func(str string, start, length int) string {
			if start >= len(str) {
				return ""
			}
			end := start + length
			if end > len(str) {
				end = len(str)
			}
			return str[start:end]
		},
		"formatDate": func(t interface{}) string {
			// 这里可以添加日期格式化逻辑
			return fmt.Sprintf("%v", t)
		},
	}
}

// 核心修改：loadTemplates（支持多级目录，带路径校验和日志）
func (tm *TemplateManager) LoadTemplates(r *gin.Engine, funcs template.FuncMap) {
	// 1. 先校验模板根目录是否存在
	if _, err := os.Stat(tm.templateDir); os.IsNotExist(err) {
		panic(fmt.Sprintf("模板根目录不存在：%s（当前工作目录：%s）", tm.templateDir, getCurrentDir()))
	}

	// 2. 创建模板集合并注册自定义函数
	templ := template.New("").Funcs(funcs)

	// 3. 遍历每个模板子目录（如 "default"、"simple"）
	for _, subDir := range tm.templateDirs {
		// 拼接完整的子目录路径（如 "templates/default"）
		rootPath := filepath.Join(tm.templateDir, subDir)

		// 关键日志：打印当前要遍历的目录路径
		fmt.Printf("开始遍历模板子目录：%s（是否存在？%t）\n", rootPath, isDirExist(rootPath))

		// 校验子目录是否存在（不存在则跳过或报错）
		if !isDirExist(rootPath) {
			panic(fmt.Sprintf("模板子目录不存在：%s（父目录：%s）", rootPath, tm.templateDir))
		}

		// 4. 递归遍历子目录下的所有文件
		err := filepath.Walk(rootPath, func(filePath string, info os.FileInfo, walkErr error) error {
			// 先处理遍历过程中的错误（如权限问题）
			if walkErr != nil {
				fmt.Printf("遍历文件时出错：%s（路径：%s）\n", walkErr.Error(), filePath)
				return walkErr
			}

			// 只处理「非目录」且「后缀为.html」的文件
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".html") {
				// 关键日志：打印找到的HTML文件路径
				fmt.Printf("找到HTML模板文件：%s\n", filePath)

				// 计算模板名称（相对于根目录的相对路径，避免重复）
				// 例如："templates/default/components/header.html" → "default/components/header.html"
				relPath, err := filepath.Rel(tm.templateDir, filePath)
				if err != nil {
					fmt.Printf("计算相对路径失败：%s（文件路径：%s）\n", err.Error(), filePath)
					return err
				}

				// 统一路径分隔符（Windows的“\”转为“/”，兼容模板名称）
				tplName := filepath.ToSlash(relPath)

				// 解析模板并添加到集合
				_, err = templ.New(tplName).ParseFiles(filePath)
				if err != nil {
					fmt.Printf("解析模板失败：%s（模板文件：%s，模板名称：%s）\n", err.Error(), filePath, tplName)
					return err
				}
			}
			return nil
		})

		// 处理遍历目录的整体错误
		if err != nil {
			panic(fmt.Sprintf("遍历模板子目录失败：%s（目录：%s）\n", err.Error(), rootPath))
		}
	}

	// 5. 将加载好的模板集合设置到Gin引擎
	r.SetHTMLTemplate(templ)
	fmt.Printf("模板加载完成！共加载 %d 个模板文件\n", len(templ.Templates()))
}

// 辅助函数：判断目录是否存在
func isDirExist(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// 辅助函数：获取当前工作目录（帮你确认程序执行路径）
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "获取工作目录失败：" + err.Error()
	}
	return dir
}

// 检查模板是否有效
func (tm *TemplateManager) isValidTemplate(name string) bool {
	for _, dir := range tm.templateDirs {
		if dir == name {
			return true
		}
	}
	return false
}
