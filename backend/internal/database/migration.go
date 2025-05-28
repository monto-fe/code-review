package database

import (
	"log"

	"code-review-go/internal/model"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	log.Println("Starting database migration...")

	// 在这里添加需要迁移的模型
	err := DB.AutoMigrate(
		&model.User{},
		// 在这里添加其他模型
	)

	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}
