package user

import (
	"strconv"

	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/constants"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/pkg/utils"
	"code-review-go/internal/service"

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
		response.Error(c, err, "Invalid request parameters", int(constants.RetCodeInvalidParams))
		return
	}
	userService := service.NewUserService(database.DB)
	findData, err := userService.FindUserByUsername(loginReq.User, loginReq.Namespace)
	hashedPassword := utils.HashPassword(loginReq.Password)
	if err != nil || findData.Password != hashedPassword {
		response.Error(c, err, "Invalid username or password", int(constants.RetCodeInvalidCredentials))
		return
	}

	loginParams := model.LoginParams{
		ID:        findData.ID,
		User:      loginReq.User,
		Namespace: loginReq.Namespace,
	}
	jwtToken, err := userService.Login(loginParams)

	if err != nil || jwtToken["jwtToken"] == "" {
		response.Error(c, err, "Login failed", int(constants.RetCodeLoginFailed))
		return
	}

	// 返回登录成功信息
	response.Success(c, gin.H{
		"jwt_token": jwtToken["jwtToken"],
		"user":      loginReq.User,
	}, "Login successful", int(constants.RetCodeSuccess))
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
		response.Error(c, nil, "Namespace is required", int(constants.RetCodeInvalidParams))
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
		response.Error(c, err, "Failed to get user list", int(constants.RetCodeGetUserListFailed))
		return
	}

	response.Success(c, gin.H{
		"data":  userList,
		"count": total,
	}, "Success", int(constants.RetCodeSuccess))
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
		response.Error(c, nil, "User ID not found in request header", int(constants.RetCodeOperatorMissing))
		return
	}

	// 转换用户ID为uint
	userIdUint, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		response.Error(c, err, "Invalid user ID format", int(constants.RetCodeInvalidUserID))
		return
	}

	// 获取用户信息和角色列表
	result, err := userService.FindUserAndRoleByID(uint(userIdUint))
	if err != nil {
		response.Error(c, err, "Failed to get user info", int(constants.RetCodeGetUserInfoFailed))
		return
	}

	response.Success(c, result, "Success", int(constants.RetCodeSuccess))
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
		response.Error(c, err, "Invalid request parameters", int(constants.RetCodeInvalidParams))
		return
	}

	// 参数格式校验
	if err := validateCreateInnerUserRequest(&req); err != nil {
		response.Error(c, err, "参数格式校验失败", int(constants.RetCodeInvalidParams))
		return
	}

	operator := c.GetHeader("remoteUser")
	if operator == "" {
		response.Error(c, nil, "Operator not found in request header", int(constants.RetCodeOperatorMissing))
		return
	}

	userService := service.NewUserService(database.DB)

	// 检查用户是否已存在
	exists, err := userService.CheckUsernameExists(req.Namespace, req.User)
	if err != nil {
		response.Error(c, err, "Failed to check user existence", int(constants.RetCodeCheckUserExistenceFailed))
		return
	}
	if exists {
		response.Error(c, nil, "User already exists", int(constants.RetCodeUserAlreadyExists))
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
		response.Error(c, err, "Failed to create user", int(constants.RetCodeUserOperationFailed))
		return
	}

	response.Success(c, result, "User created successfully", int(constants.RetCodeSuccess))
}

// validateCreateInnerUserRequest 验证创建内部用户请求参数
func validateCreateInnerUserRequest(req *dto.CreateInnerUserRequest) error {
	// 验证命名空间格式
	if err := utils.ValidateNamespace(req.Namespace); err != nil {
		return err
	}

	// 验证用户名格式
	if err := utils.ValidateUsername(req.User); err != nil {
		return err
	}

	// 验证姓名格式
	if err := utils.ValidateName(req.Name); err != nil {
		return err
	}

	// 验证职位格式
	if err := utils.ValidateJob(req.Job); err != nil {
		return err
	}

	// 验证密码强度
	if err := utils.ValidatePassword(req.Password); err != nil {
		return err
	}

	// 验证邮箱格式
	if err := utils.ValidateEmail(req.Email); err != nil {
		return err
	}

	// 验证手机号格式
	if err := utils.ValidatePhone(req.PhoneNumber); err != nil {
		return err
	}

	// 验证角色ID数组
	if err := utils.ValidateRoleIDs(req.RoleIDs); err != nil {
		return err
	}

	return nil
}

// validateUpdateInnerUserRequest 验证更新内部用户请求参数
func validateUpdateInnerUserRequest(req *dto.UpdateInnerUserRequest) error {
	// 验证命名空间格式
	if err := utils.ValidateNamespace(req.Namespace); err != nil {
		return err
	}

	// 验证用户名格式（如果提供）
	if req.User != "" {
		if err := utils.ValidateUsername(req.User); err != nil {
			return err
		}
	}

	// 验证姓名格式（如果提供）
	if req.Name != "" {
		if err := utils.ValidateName(req.Name); err != nil {
			return err
		}
	}

	// 验证职位格式（如果提供）
	if req.Job != "" {
		if err := utils.ValidateJob(req.Job); err != nil {
			return err
		}
	}

	// 验证密码强度（如果提供）
	if req.Password != nil && *req.Password != "" {
		if err := utils.ValidatePassword(*req.Password); err != nil {
			return err
		}
	}

	// 验证邮箱格式（如果提供）
	if req.Email != nil && *req.Email != "" {
		if err := utils.ValidateEmail(*req.Email); err != nil {
			return err
		}
	}

	// 验证手机号格式（如果提供）
	if req.PhoneNumber != nil && *req.PhoneNumber != "" {
		if err := utils.ValidatePhone(*req.PhoneNumber); err != nil {
			return err
		}
	}

	// 验证角色ID数组
	if err := utils.ValidateRoleIDs(req.RoleIDs); err != nil {
		return err
	}

	return nil
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
		response.Error(c, err, "Invalid request parameters", int(constants.RetCodeInvalidParams))
		return
	}

	// 参数格式校验
	if err := validateUpdateInnerUserRequest(&req); err != nil {
		response.Error(c, err, "参数格式校验失败", int(constants.RetCodeInvalidParams))
		return
	}

	operator := c.GetHeader("remoteUser")
	if operator == "" {
		response.Error(c, nil, "Operator not found in request header", int(constants.RetCodeOperatorMissing))
		return
	}

	userService := service.NewUserService(database.DB)

	// 检查用户是否存在
	exists, err := userService.CheckUsernameExists(req.Namespace, req.User)
	if err != nil {
		response.Error(c, err, "Failed to check user existence", int(constants.RetCodeCheckUserExistenceFailed))
		return
	}
	if !exists {
		response.Error(c, nil, "User not found", int(constants.RetCodeUserAlreadyExists))
		return
	}

	// 更新用户
	updateReq := model.UpdateUserReq{
		ID:          req.ID,
		User:        req.User,
		Namespace:   req.Namespace,
		Name:        req.Name,
		Job:         req.Job,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		RoleIDs:     req.RoleIDs,
		Operator:    operator,
	}

	// 只有当密码字段不为空时才更新密码
	if req.Password != nil && *req.Password != "" {
		hashedPassword := utils.HashPassword(*req.Password)
		updateReq.Password = &hashedPassword
	}

	if err := userService.UpdateInnerUser(updateReq); err != nil {
		response.Error(c, err, "Failed to update user", int(constants.RetCodeUserOperationFailed))
		return
	}

	response.Success(c, nil, "User updated successfully", int(constants.RetCodeSuccess))
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
		response.Error(c, err, "参数错误", int(constants.RetCodeBadRequest))
		return
	}

	userService := service.NewUserService(database.DB)
	if err := userService.DeleteUser(req.ID, req.Namespace, req.User); err != nil {
		response.Error(c, err, "Failed to delete user", int(constants.RetCodeDeleteUserFailed))
		return
	}

	response.Success(c, nil, "User deleted successfully", int(constants.RetCodeSuccess))
}
