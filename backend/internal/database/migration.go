package database

import (
	"fmt"
	"log"

	"code-review-go/internal/model"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	log.Println("Starting database migration...")

	// 在这里添加需要迁移的模型
	err := DB.AutoMigrate(
		&model.User{},
		&model.AImessage{}, // 添加AIMessage模型
		// 在这里添加其他模型
	)

	if err != nil {
		return err
	}

	// 创建AIMessage表的索引以优化查询性能
	if err := createAIMessageIndexes(DB); err != nil {
		log.Printf("Warning: Failed to create AIMessage indexes: %v", err)
		// 不返回错误，因为索引创建失败不应该阻止应用启动
	}

	log.Println("Database migration completed successfully")
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
