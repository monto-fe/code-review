package user

import (
	"code-review-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户相关路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 登录路由
	r.POST("/login", Login)

	// 用户管理路由
	users := r.Group("/user")
	{
		// 需要认证的路由
		users.Use(middleware.AuthenticateJWT())
		{
			users.GET("", GetUserList)      // 获取用户列表
			users.POST("", CreateInnerUser) // 创建内部用户
			users.PUT("", UpdateInnerUser)  // 更新内部用户
			users.DELETE("", DeleteUser)    // 删除用户
			users.GET("/info", GetUserInfo) // 获取用户信息
		}
	}
	// role
	role := r.Group("/role")
	{
		role.GET("", GetRoleList)   // 获取角色列表
		role.POST("", CreateRole)   // 创建角色
		role.PUT("", UpdateRole)    // 更新角色
		role.DELETE("", DeleteRole) // 删除角色
	}

	// resource
	resource := r.Group("/resource")
	{
		resource.GET("", GetResourceList)   // 获取资源列表
		resource.POST("", CreateResource)   // 创建资源
		resource.PUT("", UpdateResource)    // 更新资源
		resource.DELETE("", DeleteResource) // 删除资源
	}
}
