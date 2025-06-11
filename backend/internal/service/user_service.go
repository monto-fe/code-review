package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"code-review-go/internal/dto"
	"code-review-go/internal/middleware"
	"code-review-go/internal/model"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// generateSignToken 生成JWT token
func generateSignToken(params model.LoginParams, now int64) (string, error) {
	claims := jwt.MapClaims{
		"id":          params.ID,
		"user":        params.User,
		"namespace":   params.Namespace,
		"create_time": now,
		"exp":         now + 8*60*60, // 8小时过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.TokenSecretKey)
}

// FindUserByID 根据ID查找用户
func (s *UserService) FindUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

// FindUserAndRoleByID 查找用户及其角色
func (s *UserService) FindUserAndRoleByID(id uint) (map[string]interface{}, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	var userRoles []model.UserRole
	if err := s.db.Where("user = ?", user.User).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	var roles []model.Role
	if err := s.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"userInfo": user,
		"roleList": roles,
	}, nil
}

// FindUserByUsername 根据用户名查找用户
func (s *UserService) FindUserByUsername(user, namespace string) (*model.User, error) {
	var userModel model.User
	err := s.db.Where("`user` = ? AND namespace = ?", user, namespace).First(&userModel).Error
	if err != nil {
		return nil, err
	}
	return &userModel, nil
}

// CheckUsernameExists 检查用户名是否存在
func (s *UserService) CheckUsernameExists(namespace, username string) (bool, error) {
	var count int64
	err := s.db.Model(&model.User{}).Where("namespace = ? AND user = ?", namespace, username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckEmailExists 检查邮箱是否存在
func (s *UserService) CheckEmailExists(email string) (bool, error) {
	var count int64
	err := s.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// Login 用户登录
func (s *UserService) Login(params model.LoginParams) (map[string]string, error) {
	var token model.Token
	err := s.db.Where("user = ?", params.User).Order("create_time DESC").First(&token).Error

	now := time.Now().Unix()
	if err == nil && token.ExpiredAt > now {
		return map[string]string{"jwtToken": token.Token}, nil
	}

	// 生成新token
	jwtToken, err := generateSignToken(params, now)
	if err != nil {
		return map[string]string{"jwtToken": ""}, nil
	}

	// 创建新token记录
	newToken := model.Token{
		Token:      jwtToken,
		User:       params.User,
		ExpiredAt:  now + 8*60*60, // 8小时过期
		CreateTime: now,
		UpdateTime: now,
	}

	if err := s.db.Create(&newToken).Error; err != nil {
		return map[string]string{"jwtToken": ""}, nil
	}

	return map[string]string{"jwtToken": jwtToken}, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *model.User) error {
	return s.db.Create(user).Error
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(req model.UpdateUserReq) error {
	return s.db.Model(&model.User{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"user":         req.User,
		"name":         req.Name,
		"email":        req.Email,
		"phone_number": req.PhoneNumber,
		"job":          req.Job,
	}).Error
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(query model.UserListQuery) ([]dto.UserListItem, int64, error) {
	var users []model.User
	var count int64
	fmt.Println("query", query)

	db := s.db.Model(&model.User{}).Where("namespace = ?", query.Namespace)
	if query.User != "" {
		db = db.Where("user = ?", query.User)
	}
	if query.UserName != "" {
		db = db.Where("name LIKE ?", "%"+query.UserName+"%")
	}
	if query.RoleName != "" {
		// 通过角色名称关联查询
		db = db.Joins("JOIN t_user_role ON t_user.user = t_user_role.user AND t_user.namespace = t_user_role.namespace").
			Joins("JOIN t_role ON t_user_role.role_id = t_role.id").
			Where("t_role.name LIKE ?", "%"+query.RoleName+"%")
	}

	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((query.Current - 1) * query.PageSize).
		Limit(query.PageSize).
		Order("id DESC").
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为不含密码的DTO
	var userList []dto.UserListItem
	for _, u := range users {
		userList = append(userList, dto.UserListItem{
			ID:          u.ID,
			User:        u.User,
			Name:        u.Name,
			Job:         u.Job,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			Namespace:   u.Namespace,
		})
	}

	return userList, count, nil
}

// CreateInnerUser 创建内部用户
func (s *UserService) CreateInnerUser(params model.CreateUserParams) (interface{}, error) {
	// 开启事务
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	now := time.Now().Unix()

	// 创建用户
	user := model.User{
		Namespace:   params.Namespace,
		User:        params.User,
		Name:        params.Name,
		Job:         params.Job,
		Password:    params.Password,
		Email:       params.Email,
		PhoneNumber: params.PhoneNumber,
		CreateTime:  now,
		UpdateTime:  now,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果有角色ID，创建用户角色关联
	if len(params.RoleIDs) > 0 {
		for _, roleID := range params.RoleIDs {
			userRole := model.UserRole{
				Namespace:  params.Namespace,
				User:       params.User,
				RoleID:     roleID,
				Status:     1, // 启用状态
				Operator:   params.Operator,
				CreateTime: now,
				UpdateTime: now,
			}
			if err := tx.Create(&userRole).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateInnerUser 更新内部用户
func (s *UserService) UpdateInnerUser(req model.UpdateUserReq) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新用户信息
		updates := map[string]interface{}{
			"password":    req.Password,
			"name":        req.Name,
			"job":         req.Job,
			"update_time": time.Now().Unix(),
		}
		if req.Email != nil {
			updates["email"] = *req.Email
		}
		if req.PhoneNumber != nil {
			updates["phone_number"] = *req.PhoneNumber
		}

		if err := tx.Model(&model.User{}).Where("id = ? AND namespace = ?", req.ID, req.Namespace).
			Updates(updates).Error; err != nil {
			return err
		}

		// 更新角色信息
		if len(req.RoleIDs) > 0 {
			// 删除旧的角色关联
			if err := tx.Where("namespace = ? AND user = ?", req.Namespace, req.User).
				Delete(&model.UserRole{}).Error; err != nil {
				return err
			}

			// 创建新的角色关联
			now := time.Now().Unix()
			for _, roleID := range req.RoleIDs {
				userRole := model.UserRole{
					Namespace:  req.Namespace,
					User:       req.User,
					RoleID:     roleID,
					Status:     1, // 启用状态
					Operator:   req.Operator,
					CreateTime: now,
					UpdateTime: now,
				}
				if err := tx.Create(&userRole).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint, namespace, user string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 先检查用户是否存在
		var count int64
		if err := tx.Model(&model.User{}).Where("id = ? AND namespace = ? AND user = ?", id, namespace, user).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("user not found")
		}

		// 删除用户
		if err := tx.Where("id = ? AND namespace = ? AND user = ?", id, namespace, user).
			Delete(&model.User{}).Error; err != nil {
			return err
		}

		// 删除用户角色关联
		if err := tx.Where("namespace = ? AND user = ?", namespace, user).
			Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindUserByUserName 根据用户名查找用户
func (s *UserService) FindUserByUserName(userName string) ([]string, error) {
	var users []string
	err := s.db.Model(&model.User{}).Where("name LIKE ?", "%"+userName+"%").
		Pluck("user", &users).Error
	return users, err
}
