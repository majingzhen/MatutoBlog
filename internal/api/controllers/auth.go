package controllers

import (
	"matuto-blog/internal/api/middlewares"
	"matuto-blog/internal/database"
	"matuto-blog/internal/models"
	"matuto-blog/pkg/common"
	"matuto-blog/pkg/utils"

	"github.com/gin-gonic/gin"
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
		common.ServerError(c, "参数错误: "+err.Error())
		return
	}

	// 查找用户
	var user models.User
	if err := database.DB.Where("account = ?", req.Account).First(&user).Error; err != nil {
		common.ServerError(c, "账户名或密码错误")
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		common.ServerError(c, "账户名或密码错误")
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		common.ServerError(c, "账户已被禁用")
		return
	}

	// 生成token
	token, err := middlewares.GenerateToken(&user)
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

// AdminLogout 管理员退出登录
func (a *AuthController) Logout(c *gin.Context) {
	common.SuccessWithMessage(c, "退出登录成功", nil)
}

// GetProfile 获取用户信息
func (a *AuthController) GetProfile(c *gin.Context) {
	value, exists := c.Get("user")
	if !exists {
		common.ServerError(c, "未找到用户信息")
		return
	}

	user, ok := value.(*models.User)
	if !ok {
		common.ServerError(c, "未找到用户信息")
		return
	}
	common.SuccessWithMessage(c, "获取成功", gin.H{
		"id":       user.Id,
		"username": user.Username,
		"account":  user.Account,
		"email":    user.Email,
		"avatar":   user.Avatar,
	})
}
