package database

import (
	"fmt"
	"log"

	"code-review-go/internal/model"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	log.Println("开始数据库迁移...")

	// 需要迁移的模型列表
	models := []interface{}{
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.Resource{},
		&model.Token{},
		&model.GitlabInfo{},
		&model.AIConfig{},
		&model.AIManager{},
		&model.AIMessage{},
		&model.CustomRule{},
		&model.CommonRule{},
		&model.Namespace{},
		&model.ManualCheckTask{}, // 新增手动审核任务表
		&model.RolePermission{},  // 角色权限关联表
	}

	// 执行迁移
	for _, m := range models {
		if err := DB.AutoMigrate(m); err != nil {
			return fmt.Errorf("迁移表 %T 失败: %w", m, err)
		}
		log.Printf("成功迁移表: %T", m)
	}

	// 创建AIMessage表的索引以优化查询性能
	if err := createAIMessageIndexes(DB); err != nil {
		log.Printf("Warning: Failed to create AIMessage indexes: %v", err)
		// 不返回错误，因为索引创建失败不应该阻止应用启动
	}

	log.Println("数据库迁移完成")
	return nil
}

// 添加新的索引创建函数
func createAIMessageIndexes(db *gorm.DB) error {
	// 为AIMessage表创建复合索引，优化项目命名空间查询
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_ai_message_project_namespace 
		ON t_ai_message (project_id, project_namespace)
	`).Error; err != nil {
		return fmt.Errorf("创建AIMessage项目命名空间索引失败: %v", err)
	}

	// 为AIMessage表创建项目ID索引
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_ai_message_project_id 
		ON t_ai_message (project_id)
	`).Error; err != nil {
		return fmt.Errorf("创建AIMessage项目ID索引失败: %v", err)
	}

	// 为AIMessage表创建创建时间索引，优化时间范围查询
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_ai_message_create_time 
		ON t_ai_message (create_time)
	`).Error; err != nil {
		return fmt.Errorf("创建AIMessage创建时间索引失败: %v", err)
	}

	return nil
}
