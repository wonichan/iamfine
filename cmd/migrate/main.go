package main

import (
	"context"
	"fmt"
	"hupu/services/post/repository"
	"hupu/shared/utils"
	"log"
)

func main() {
	// 初始化数据库连接
	if err := utils.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建迁移器
	migrator := repository.NewDatabaseMigrator()
	ctx := context.Background()

	fmt.Println("🚀 开始数据库索引优化...")

	// 执行所有迁移
	if err := migrator.RunAllMigrations(ctx); err != nil {
		log.Fatalf("❌ 数据库迁移失败: %v", err)
	}

	fmt.Println("🎉 数据库索引优化完成！")

	// 输出索引创建状态报告
	printIndexReport()
}

func printIndexReport() {
	fmt.Println("\n📋 索引优化报告:")
	fmt.Println("================")

	fmt.Println("\n✅ GORM自动创建的索引:")
	gormIndexes := []string{
		"idx_user_created (user_id, created_at) - 用户帖子查询",
		"idx_category_created (category, created_at) - 分类帖子查询",
		"idx_anonymous_created (is_anonymous, created_at) - 匿名帖子查询",
		"idx_category_hot (category, like_count, comment_count, created_at) - 热门帖子",
		"idx_post_score (post_id, score) - 评分统计",
		"idx_user_post_rating (user_id, post_id) - 用户评分唯一约束",
		"idx_user_post_favorite (user_id, post_id) - 用户收藏唯一约束",
	}

	for _, index := range gormIndexes {
		fmt.Printf("  • %s\n", index)
	}

	fmt.Println("\n🔧 手动创建的高级索引:")
	advancedIndexes := []string{
		"idx_posts_title_content - 全文搜索索引",
		"idx_posts_tags_generated - JSON标签搜索索引",
		"idx_posts_hot_score - 热门度计算索引",
		"idx_posts_active - 活跃帖子条件索引",
		"idx_posts_today - 今日帖子索引",
		"idx_posts_anonymous_active - 活跃匿名帖子索引",
	}

	for _, index := range advancedIndexes {
		fmt.Printf("  • %s\n", index)
	}

	fmt.Println("\n💡 使用建议:")
	suggestions := []string{
		"定期执行 ANALYZE TABLE 更新统计信息",
		"监控慢查询日志，根据实际使用情况调整索引",
		"在生产环境中逐步创建索引，避免影响业务",
		"考虑在业务低峰期执行索引优化",
	}

	for _, suggestion := range suggestions {
		fmt.Printf("  • %s\n", suggestion)
	}
}
