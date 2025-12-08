package blog

import (
	"fmt"
	model "job/blog/Model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	var err error
	db, err = gorm.Open(mysql.Open("root:123456@(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err != nil {
		return err
	}
	fmt.Println("Database migrated successfully")
	return nil
}

func GetDB() *gorm.DB {
	return db
}
