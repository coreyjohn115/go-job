package main

import (
	"fmt"
	blog "job/blog"

	"github.com/gin-gonic/gin"
)

func main() {
	err := blog.InitDB()
	if err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	engine := gin.Default()
	blog.RegisterRoutes(engine)
	err = engine.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("启动服务器失败: %v", err))
	}
}
