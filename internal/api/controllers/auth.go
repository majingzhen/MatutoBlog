package controllers

import (
	"matuto-blog/internal/api/middleware"
	"net/http"

	"matuto-blog/internal/database"
	"matuto-blog/internal/models"
	"matuto-blog/pkg/common"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthController 认证控制器
type AuthController struct{}

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func (a *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 查找用户
	var user models.User
	if err := database.DB.Where("account = ?", req.Account).First(&user).Error; err != nil {
		common.Unauthorized(c, "账户名或密码错误")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		common.Unauthorized(c, "账户名或密码错误")
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		common.Unauthorized(c, "账户已被禁用")
		return
	}

	// 生成token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		common.ServerError(c, "生成令牌失败")
		return
	}

	// 返回登录信息
	common.SuccessWithMessage(c, "登录成功", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.Id,
			"username": user.Username,
			"account":  user.Account,
			"email":    user.Email,
			"avatar":   user.Avatar,
		},
	})
}

// AdminLogin 管理员登录页面
func (a *AuthController) AdminLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login.html", gin.H{
		"title": "管理员登录",
	})
}

// AdminLogin 管理员登录处理
func (a *AuthController) AdminLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")

	if account == "" || password == "" {
		c.HTML(http.StatusBadRequest, "admin/login.html", gin.H{
			"title": "管理员登录",
			"error": "用户名和密码不能为空",
		})
		return
	}

	// 查找用户
	var user models.User
	if err := database.DB.Where("account = ?", account).First(&user).Error; err != nil {
		c.HTML(http.StatusBadRequest, "admin/login.html", gin.H{
			"title": "管理员登录",
			"error": "账户名或密码错误",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusBadRequest, "admin/login.html", gin.H{
			"title": "管理员登录",
			"error": "账户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.HTML(http.StatusBadRequest, "admin/login.html", gin.H{
			"title": "管理员登录",
			"error": "账户已被禁用",
		})
		return
	}

	// 生成token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin/login.html", gin.H{
			"title": "管理员登录",
			"error": "系统错误，请稍后再试",
		})
		return
	}

	// 设置cookie
	c.SetCookie("admin_session", token, 86400*7, "/", "", false, true) // 7天有效期

	// 重定向到管理后台首页
	c.Redirect(http.StatusFound, "/admin")
}

// AdminLogout 管理员退出登录
func (a *AuthController) AdminLogout(c *gin.Context) {
	// 清除cookie
	c.SetCookie("admin_session", "", -1, "/", "", false, true)

	// 重定向到登录页面
	c.Redirect(http.StatusFound, "/admin/login")
}

// GetProfile 获取用户信息
func (a *AuthController) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		common.Unauthorized(c, "未登录")
		return
	}

	u := user.(*models.User)
	common.SuccessWithMessage(c, "获取成功", gin.H{
		"id":       u.Id,
		"username": u.Username,
		"account":  u.Account,
		"email":    u.Email,
		"avatar":   u.Avatar,
	})
}
