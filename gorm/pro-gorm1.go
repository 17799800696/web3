package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/test/init_project/config"
)

// User model definition
type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"uniqueIndex;size:255"`
	Password string
	Posts    []Post `gorm:"foreignKey:UserID"`
}

// Post model definition
type Post struct {
	gorm.Model
	Title         string `gorm:"size:255"`
	Content       string `gorm:"type:text"`
	UserID        uint
	User          User
	Comments      []Comment `gorm:"foreignKey:PostID"`
	CommentCount  int       `gorm:"default:0;check:comment_count >= 0"`                   // add constraint to ensure the count is not negative
	CommentStatus string    `gorm:"default:'有评论';check:comment_status IN ('有评论', '无评论')"` // add constraint to ensure the status is valid
}

// Comment model definition
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text"`
	UserID  uint
	PostID  uint
	User    User `gorm:"foreignKey:UserID"`
	Post    Post `gorm:"foreignKey:PostID"`
}

// safer comment create hook function
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	// use atomic operation to increase the comment count, avoid race condition
	result := tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))

	if result.Error != nil {
		return fmt.Errorf("update post comment count failed: %v", result.Error)
	}

	// check the updated rows, ensure the post exists
	if result.RowsAffected == 0 {
		return fmt.Errorf("post with id %d not found", c.PostID)
	}

	// use transaction level query to check the status
	var post Post
	if err := tx.Select("comment_count, comment_status").
		Where("id = ?", c.PostID).
		First(&post).Error; err != nil {
		return fmt.Errorf("query post failed: %v", err)
	}

	// if the comment count is 1, update the status
	if post.CommentCount == 1 && post.CommentStatus != "有评论" {
		if err := tx.Model(&Post{}).
			Where("id = ? AND comment_status != ?", c.PostID, "有评论").
			Update("comment_status", "有评论").Error; err != nil {
			return fmt.Errorf("update post status failed: %v", err)
		}
		fmt.Printf("post id: %d received the first comment, status updated to '有评论'\n", c.PostID)
	}

	return nil
}

// safer comment delete hook function
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// use atomic operation to reduce the comment count, avoid race condition
	result := tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))

	if result.Error != nil {
		return fmt.Errorf("update post comment count failed: %v", result.Error)
	}

	// check the updated rows, ensure the post exists
	if result.RowsAffected == 0 {
		return fmt.Errorf("post with id %d not found", c.PostID)
	}

	// use transaction level query to check the status
	var post Post
	if err := tx.Select("comment_count, comment_status").
		Where("id = ?", c.PostID).
		First(&post).Error; err != nil {
		return fmt.Errorf("query post failed: %v", err)
	}

	// if the comment count is 0, update the status
	if post.CommentCount == 0 && post.CommentStatus != "无评论" {
		if err := tx.Model(&Post{}).
			Where("id = ? AND comment_status != ?", c.PostID, "无评论").
			Update("comment_status", "无评论").Error; err != nil {
			return fmt.Errorf("update post status failed: %v", err)
		}
		fmt.Printf("post id: %d comment is empty, status updated to '无评论'\n", c.PostID)
	}

	return nil
}

func printUserPostsWithComments(user User) {
	fmt.Printf("user %s's posts:\n", user.Name)
	for _, post := range user.Posts {
		fmt.Printf("- %s (comment count: %d):\n", post.Title, len(post.Comments))
		for _, comment := range post.Comments {
			fmt.Printf("  - %s (user id: %d)\n", comment.Content, comment.UserID)
		}
	}
}

func insertTestData(db *gorm.DB) error {
	// check if there is data
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount > 0 {
		return nil // there is data, not insert again
	}

	// create users
	users := []User{
		{Name: "张三", Email: "zhangsan@example.com", Password: "123456"},
		{Name: "李四", Email: "lisi@example.com", Password: "123456"},
		{Name: "王五", Email: "wangwu@example.com", Password: "123456"},
	}
	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	// create posts
	posts := []Post{
		{Title: "GORM 入门教程", Content: "这是一篇关于GORM的入门教程...", UserID: 1, CommentCount: 2},
		{Title: "Go语言实战", Content: "介绍Go语言的高级特性和最佳实践...", UserID: 1, CommentCount: 1},
		{Title: "微服务架构", Content: "讨论微服务设计模式和实现方法...", UserID: 2, CommentCount: 3},
	}
	for _, post := range posts {
		if err := db.Create(&post).Error; err != nil {
			return err
		}
	}

	// create comments
	comments := []Comment{
		{Content: "非常实用的教程！", UserID: 2, PostID: 1},
		{Content: "期待后续更新！", UserID: 3, PostID: 1},
		{Content: "适合中级开发者", UserID: 2, PostID: 2},
		{Content: "架构设计很清晰", UserID: 1, PostID: 3},
		{Content: "需要更多代码示例", UserID: 3, PostID: 3},
		{Content: "微服务确实复杂", UserID: 2, PostID: 3},
	}
	for _, comment := range comments {
		if err := db.Create(&comment).Error; err != nil {
			return err
		}
	}

	fmt.Println("insert test data successfully")
	return nil
}

func main() {
	// use global config to get database connection
	db := config.GetDB()

	// auto migrate model to database
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal("table structure migration failed:", err)
	}

	log.Println("table created successfully")

	// insert test data
	if err := insertTestData(db); err != nil {
		log.Fatal("insert test data failed:", err)
	}

	var user User
	if err := db.Preload("Posts.Comments").First(&user, 1).Error; err != nil {
		log.Fatal("user not found:", err)
	}
	printUserPostsWithComments(user)

	// query the most commented article
	var topPost Post
	if err := db.Model(&Post{}).
		Order("comment_count DESC").
		First(&topPost).Error; err != nil {
		log.Fatal("query failed:", err)
	}
	fmt.Printf("\nthe most commented article: %s (comment count: %d)\n", topPost.Title, topPost.CommentCount)
}
