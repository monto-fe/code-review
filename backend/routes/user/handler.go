package user

import (
	"strconv"

	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
	"code-review-go/internal/service"

	"code-review-go/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param login body dto.UserLoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "返回 JWT token 和用户信息"
// @Router /v1/login [post]
func Login(c *gin.Context) {
	var loginReq dto.UserLoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		response.Error(c, err, "Invalid request parameters", 10001)
		return
	}
	userService := service.NewUserService(database.DB)
	findData, err := userService.FindUserByUsername(loginReq.User, loginReq.Namespace)
	hashedPassword := utils.HashPassword(loginReq.Password)
	if err != nil || findData.Password != hashedPassword {
		response.Error(c, err, "Invalid username or password", 10002)
		return
	}

	loginParams := model.LoginParams{
		ID:        findData.ID,
		User:      loginReq.User,
		Namespace: loginReq.Namespace,
	}
	jwtToken, err := userService.Login(loginParams)

	if err != nil || jwtToken["jwtToken"] == "" {
		response.Error(c, err, "Login failed", 10003)
		return
	}

	// 返回登录成功信息
	response.Success(c, gin.H{
		"jwt_token": jwtToken["jwtToken"],
		"user":      loginReq.User,
	}, "Login successful", 0)
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 获取指定命名空间的用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Security ApiKeyAuth
func GetUserList(c *gin.Context) {
	namespace := c.Query("namespace")
	if namespace == "" {
		response.Error(c, nil, "Namespace is required", 10001)
		return
	}

	userService := service.NewUserService(database.DB)
	query := model.UserListQuery{
		Current:   1,
		PageSize:  10,
		Namespace: namespace,
	}
	userList, total, err := userService.GetUserList(query)
	if err != nil {
		response.Error(c, err, "Failed to get user list", 10004)
		return
	}

	response.Success(c, gin.H{
		"data":  userList,
		"count": total,
	}, "Success", 0)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "返回用户信息和角色列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /v1/user/info [get]
func GetUserInfo(c *gin.Context) {
	userService := service.NewUserService(database.DB)

	// 从请求头获取用户ID
	userId := c.GetHeader("userId")
	if userId == "" {
		response.Error(c, nil, "User ID not found in request header", 10006)
		return
	}

	// 转换用户ID为uint
	userIdUint, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		response.Error(c, err, "Invalid user ID format", 10007)
		return
	}

	// 获取用户信息和角色列表
	result, err := userService.FindUserAndRoleByID(uint(userIdUint))
	if err != nil {
		response.Error(c, err, "Failed to get user info", 10005)
		return
	}

	response.Success(c, result, "Success", 0)
}

// CreateInnerUser 创建内部用户
// @Summary 创建内部用户
// @Description 创建新的内部用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Security ApiKeyAuth
// @Param user body dto.CreateInnerUserRequest true "用户信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /v1/user [post]
func CreateInnerUser(c *gin.Context) {
	var req dto.CreateInnerUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "Invalid request parameters", 10001)
		return
	}

	operator := c.GetHeader("remoteUser")
	if operator == "" {
		response.Error(c, nil, "Operator not found in request header", 10006)
		return
	}

	userService := service.NewUserService(database.DB)

	// 检查用户是否已存在
	exists, err := userService.CheckUsernameExists(req.Namespace, req.User)
	if err != nil {
		response.Error(c, err, "Failed to check user existence", 10009)
		return
	}
	if exists {
		response.Error(c, nil, "User already exists", 10010)
		return
	}

	// 创建用户
	createParams := model.CreateUserParams{
		Namespace:   req.Namespace,
		User:        req.User,
		Name:        req.Name,
		Job:         req.Job,
		Password:    utils.HashPassword(req.Password),
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		RoleIDs:     req.RoleIDs,
		Operator:    operator,
	}

	result, err := userService.CreateInnerUser(createParams)
	if err != nil {
		response.Error(c, err, "Failed to create user", 10011)
		return
	}

	response.Success(c, result, "User created successfully", 0)
}

// UpdateInnerUser 更新内部用户
// @Summary 更新内部用户
// @Description 更新内部用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Security ApiKeyAuth
// @Param user body dto.UpdateInnerUserRequest true "用户信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /v1/user [put]
func UpdateInnerUser(c *gin.Context) {
	var req dto.UpdateInnerUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "Invalid request parameters", 10001)
		return
	}

	operator := c.GetHeader("remoteUser")
	if operator == "" {
		response.Error(c, nil, "Operator not found in request header", 10006)
		return
	}

	userService := service.NewUserService(database.DB)

	// 检查用户是否存在
	exists, err := userService.CheckUsernameExists(req.Namespace, req.User)
	if err != nil {
		response.Error(c, err, "Failed to check user existence", 10009)
		return
	}
	if !exists {
		response.Error(c, nil, "User not found", 10010)
		return
	}

	// 更新用户
	updateReq := model.UpdateUserReq{
		ID:          req.ID,
		User:        req.User,
		Namespace:   req.Namespace,
		Name:        req.Name,
		Job:         req.Job,
		Password:    utils.HashPassword(req.Password),
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		RoleIDs:     req.RoleIDs,
		Operator:    operator,
	}

	if err := userService.UpdateInnerUser(updateReq); err != nil {
		response.Error(c, err, "Failed to update user", 10011)
		return
	}

	response.Success(c, nil, "User updated successfully", 0)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Security ApiKeyAuth
// @Param body body dto.DeleteUserRequest true "删除用户参数"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /v1/user [delete]
func DeleteUser(c *gin.Context) {
	var req dto.DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}

	userService := service.NewUserService(database.DB)
	if err := userService.DeleteUser(req.ID, req.Namespace, req.User); err != nil {
		response.Error(c, err, "Failed to delete user", 10013)
		return
	}

	response.Success(c, nil, "User deleted successfully", 0)
}
