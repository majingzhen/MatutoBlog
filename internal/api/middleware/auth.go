package middleware

import (
	"matuto-blog/internal/database"
	"matuto-blog/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// Claims JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Account  string `json:"account"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   user.Id,
		Username: user.Username,
		Account:  user.Account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(viper.GetInt("jwt.access_token_ttl")))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    viper.GetString("jwt.issuer"),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("jwt.secret")))
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 解析Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "认证令牌格式错误",
			})
			c.Abort()
			return
		}

		// 解析token
		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "认证令牌无效",
			})
			c.Abort()
			return
		}

		// 验证用户是否存在
		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "用户不存在",
			})
			c.Abort()
			return
		}

		// 检查用户状态
		if user.Status != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "用户已被禁用",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user", &user)
		c.Set("user_id", user.Id)
		c.Set("username", user.Username)
		c.Set("account", user.Account)

		c.Next()
	}
}

// SessionAuth Session认证中间件 (用于管理后台)
func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从cookie或session中获取用户信息
		session, err := c.Cookie("admin_session")
		if err != nil || session == "" {
			// 如果是AJAX请求，返回JSON
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "请先登录",
				})
				c.Abort()
				return
			}

			// 重定向到登录页面
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		// 解析session token
		claims, err := ParseToken(session)
		if err != nil {
			// 清除无效的cookie
			c.SetCookie("admin_session", "", -1, "/", "", false, true)

			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "登录已过期",
				})
				c.Abort()
				return
			}

			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		// 验证用户是否存在
		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			c.SetCookie("admin_session", "", -1, "/", "", false, true)

			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "用户不存在",
				})
				c.Abort()
				return
			}

			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		// 检查用户状态
		if user.Status != 1 {
			c.SetCookie("admin_session", "", -1, "/", "", false, true)

			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "用户已被禁用",
				})
				c.Abort()
				return
			}

			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user", &user)
		c.Set("user_id", user.Id)
		c.Set("username", user.Username)
		c.Set("account", user.Account)

		c.Next()
	}
}
