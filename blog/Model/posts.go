package blog

import "gorm.io/gorm"

// 文章表模型
type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  int    `gorm:"not null"`
}
