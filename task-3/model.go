package task3

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 用户模型
type User struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
	PostCount int
	Posts     []Post `gorm:"foreignKey:UserID"`
}

// 文章模型
type Post struct {
	ID            int `gorm:"primaryKey"`
	Title         string
	Content       string
	UserID        int
	CommentsCount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Comments      []Comment `gorm:"foreignKey:PostID"`
}

func (p *Post) AfterSave(tx *gorm.DB) (err error) {
	var count int64
	tx.Model(&Post{}).Where("user_id = ?", p.UserID).Count(&count)
	tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_count", count)
	return nil
}

// BeforeDelete - 删除钩子
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var count int64
	tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count)
	if count == 0 {
		tx.Find(&Post{}).Where("id = ?", c.PostID).Update("Content", "无评论")
	}
	return nil
}

// 评论模型
type Comment struct {
	ID      int `gorm:"primaryKey"`
	Content string
	PostID  int
	UserID  int
}

func (c *Comment) String() string {
	return fmt.Sprintf("Comment{ID: %d, Content: %s, PostID: %d, UserID: %d}", c.ID, c.Content, c.PostID, c.UserID)
}

func createDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("failed to migrate database")
	}
	fmt.Println("database migrated successfully")
	return db
}

func initData() *gorm.DB {
	db := createDB()
	db.Create(&User{Name: "John", Age: 20})
	db.Create(&User{Name: "Jane", Age: 21})
	for i := range 3 {
		db.Create(&Post{Title: fmt.Sprintf("Post %d", i), Content: fmt.Sprintf("Content %d", i), UserID: 1})
	}
	for i := range 10 {
		db.Create(&Comment{Content: fmt.Sprintf("Comment %d", i), PostID: rand.Intn(3) + 1, UserID: rand.Intn(2) + 1})
	}
	return db
}

func RunGorm() {
	db := createDB()

	// 查询评论数量最多的文章信息
	var post Post
	var commentCount int64

	// 使用子查询统计评论数量，然后找出评论最多的文章
	err := db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count DESC").
		Limit(1).
		Preload("Comments").
		Find(&post).Error

	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}

	// 获取实际的评论数量
	db.Model(&Comment{}).Where("post_id = ?", post.ID).Count(&commentCount)

	fmt.Println("=== 评论数量最多的文章信息 ===")
	fmt.Printf("文章ID: %d\n", post.ID)
	fmt.Printf("标题: %s\n", post.Title)
	fmt.Printf("内容: %s\n", post.Content)
	fmt.Printf("作者ID: %d\n", post.UserID)
	fmt.Printf("评论数量: %d\n", commentCount)
	fmt.Println("\n评论列表:")
	for _, comment := range post.Comments {
		fmt.Printf("  - %s\n", comment.String())
	}
}
