package blog

import (
	model "job/blog/Model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 登录请求结构
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// 登录响应结构
type LoginResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

// register 用户注册
func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否已存在
	if err := GetDB().Where("email = ?", req.Email).First(&model.User{}).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	GetDB().Create(&user)
	token, err := GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, LoginResponse{
		Token: token,
		User:  user,
	})
}

// login 用户登录
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 生成Token
	token, err := GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成Token失败"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

// refreshToken 刷新Token
func refreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	newToken, err := RefreshToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

// profile 获取用户信息（需要认证）
func profile(c *gin.Context) {
	user, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.UserID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// 创建文章
func createPost(c *gin.Context) {
	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 自动设置作者为当前用户
	if userID, exists := GetCurrentUserID(c); exists {
		post.UserID = userID
	}

	// 保存到数据库
	if err := GetDB().Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// 获取文章详情
func getPost(c *gin.Context) {
	postID := c.Query("id")
	if postID != "" {
		// 获取特定文章
		id, err := strconv.Atoi(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		var post model.Post
		if err := GetDB().Where("id = ?", id).First(&post).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Post details", "post": []model.Post{post}})
	} else {
		// 获取所有文章
		var posts []model.Post
		GetDB().Find(&posts)
		c.JSON(http.StatusOK, gin.H{"message": "Post details", "post": posts})
	}
}

// 更新文章
func updatePost(c *gin.Context) {
	postID := c.Query("id")
	id, err := strconv.Atoi(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var post model.Post
	GetDB().Where("id = ?", id).First(&post)
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this post"})
		return
	}
	post.Title = c.Query("title")
	post.Content = c.Query("content")
	GetDB().Save(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post updated", "post": post})
}

// 删除文章
func deletePost(c *gin.Context) {
	postID := c.Query("id")
	id, err := strconv.Atoi(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var post model.Post
	GetDB().Where("id = ?", id).First(&post)
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this post"})
		return
	}
	GetDB().Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// 创建评论
func createComment(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userID, exists := GetCurrentUserID(c); exists {
		comment.UserID = userID
	}

	GetDB().Create(&comment)
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created", "code": 0})
}

// 获取评论
func getComment(c *gin.Context) {
	postID := c.Query("id")
	id, err := strconv.Atoi(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var comments []model.Comment
	GetDB().Where("post_id = ?", id).Find(&comments)
	c.JSON(http.StatusOK, gin.H{"message": "Comment details", "code": 0, "data": comments})
}

func RegisterRoutes(r *gin.Engine) {
	// 公开路由（不需要认证）
	auth := r.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/refresh", refreshToken)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/profile", profile)
		protected.POST("/create_post", createPost)
		protected.GET("/get_post", getPost)
		protected.POST("/update_post", updatePost)
		protected.POST("/delete_post", deletePost)
		protected.POST("/create_comment", createComment)
		protected.GET("/get_comment", getComment)
	}

	// 可选认证的路由（有token时提供额外信息，没有token也能访问）
	public := r.Group("/public")
	public.Use(OptionalAuthMiddleware())
	{
		public.GET("/test", func(c *gin.Context) {
			// 如果用户已认证，可以提供个性化内容
			if userID, exists := GetCurrentUserID(c); exists {
				c.JSON(http.StatusOK, gin.H{
					"message": "test for authenticated user",
					"user_id": userID,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "Public test",
				})
			}
		})
	}
}
