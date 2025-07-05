package repository

import (
	"context"
	"fmt"
	"hupu/shared/utils"

	"gorm.io/gorm"
)

// DatabaseMigrator 数据库迁移器
type DatabaseMigrator struct {
	db *gorm.DB
}

// NewDatabaseMigrator 创建数据库迁移器
func NewDatabaseMigrator() *DatabaseMigrator {
	return &DatabaseMigrator{
		db: utils.GetDB(),
	}
}

// CreateAdvancedIndexes 创建GORM无法自动创建的高级索引
func (m *DatabaseMigrator) CreateAdvancedIndexes(ctx context.Context) error {
	// 检查是否为MySQL数据库
	if m.db.Dialector.Name() != "mysql" {
		return fmt.Errorf("advanced indexes are only supported for MySQL")
	}

	// 定义需要创建的高级索引
	indexes := []string{
		// 全文搜索索引（MySQL 5.7+）
		"CREATE FULLTEXT INDEX idx_posts_title_content ON posts(title, content)",

		// JSON标签搜索的虚拟列索引（MySQL 5.7+）
		// 先创建虚拟列，再在虚拟列上创建索引
		`ALTER TABLE posts ADD COLUMN tags_generated JSON 
		 GENERATED ALWAYS AS (tags) STORED`,
		"CREATE INDEX idx_posts_tags_generated ON posts(tags_generated)",

		// 热门度计算索引（基于表达式）
		`CREATE INDEX idx_posts_hot_score ON posts(
			((like_count * 0.5 + comment_count * 0.3 + view_count * 0.1 + share_count * 0.1))
		)`,
	}

	// 执行索引创建
	for _, indexSQL := range indexes {
		// 先检查索引是否已存在
		if err := m.db.WithContext(ctx).Exec(indexSQL).Error; err != nil {
			// 如果索引已存在，忽略错误
			if !isIndexExistsError(err) {
				return fmt.Errorf("failed to create index: %s, error: %w", indexSQL, err)
			}
		}
	}

	return nil
}

// CreatePartialIndexes 创建部分索引（条件索引）
func (m *DatabaseMigrator) CreatePartialIndexes(ctx context.Context) error {
	// MySQL 8.0+ 支持函数索引
	partialIndexes := []string{
		// 活跃帖子索引（排除软删除）
		"CREATE INDEX idx_posts_active ON posts(created_at, category) WHERE deleted_at IS NULL",

		// 今日帖子索引
		"CREATE INDEX idx_posts_today ON posts(created_at) WHERE DATE(created_at) = CURDATE()",

		// 匿名帖子索引
		"CREATE INDEX idx_posts_anonymous_active ON posts(created_at) WHERE is_anonymous = 1 AND deleted_at IS NULL",
	}

	for _, indexSQL := range partialIndexes {
		if err := m.db.WithContext(ctx).Exec(indexSQL).Error; err != nil {
			if !isIndexExistsError(err) {
				return fmt.Errorf("failed to create partial index: %s, error: %w", indexSQL, err)
			}
		}
	}

	return nil
}

// OptimizeExistingIndexes 优化现有索引
func (m *DatabaseMigrator) OptimizeExistingIndexes(ctx context.Context) error {
	optimizations := []string{
		// 分析表统计信息
		"ANALYZE TABLE posts",
		"ANALYZE TABLE post_ratings",
		"ANALYZE TABLE post_favorites",

		// 优化表
		"OPTIMIZE TABLE posts",
		"OPTIMIZE TABLE post_ratings",
		"OPTIMIZE TABLE post_favorites",
	}

	for _, sql := range optimizations {
		if err := m.db.WithContext(ctx).Exec(sql).Error; err != nil {
			// 优化失败不应该中断流程，只记录警告
			fmt.Printf("Warning: failed to execute optimization: %s, error: %v\n", sql, err)
		}
	}

	return nil
}

// DropUnusedIndexes 删除不需要的索引
func (m *DatabaseMigrator) DropUnusedIndexes(ctx context.Context) error {
	// 检查并删除可能重复或不必要的索引
	dropIndexes := []string{
		// 如果有重复的单字段索引，删除它们
		"DROP INDEX IF EXISTS idx_posts_user_id ON posts",
		"DROP INDEX IF EXISTS idx_posts_category ON posts",
	}

	for _, indexSQL := range dropIndexes {
		if err := m.db.WithContext(ctx).Exec(indexSQL).Error; err != nil {
			// 索引不存在时忽略错误
			if !isIndexNotExistsError(err) {
				return fmt.Errorf("failed to drop index: %s, error: %w", indexSQL, err)
			}
		}
	}

	return nil
}

// RunAllMigrations 执行所有数据库优化
func (m *DatabaseMigrator) RunAllMigrations(ctx context.Context) error {
	steps := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"创建高级索引", m.CreateAdvancedIndexes},
		{"创建部分索引", m.CreatePartialIndexes},
		{"删除无用索引", m.DropUnusedIndexes},
		{"优化现有索引", m.OptimizeExistingIndexes},
	}

	for _, step := range steps {
		fmt.Printf("正在执行: %s...\n", step.name)
		if err := step.fn(ctx); err != nil {
			return fmt.Errorf("执行%s失败: %w", step.name, err)
		}
		fmt.Printf("✅ %s 完成\n", step.name)
	}

	return nil
}

// 辅助函数：检查是否为索引已存在错误
func isIndexExistsError(err error) bool {
	return err != nil &&
		(contains(err.Error(), "Duplicate key name") ||
			contains(err.Error(), "already exists"))
}

// 辅助函数：检查是否为索引不存在错误
func isIndexNotExistsError(err error) bool {
	return err != nil && contains(err.Error(), "doesn't exist")
}

// 辅助函数：字符串包含检查
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsInMiddle(s, substr))))
}

func containsInMiddle(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
