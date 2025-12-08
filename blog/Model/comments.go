package blog

import "gorm.io/gorm"

// 评论表模型
type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	PostID  int    `gorm:"not null"`
	UserID  int    `gorm:"not null"`
}
