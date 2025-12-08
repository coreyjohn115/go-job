# Go Job 项目

一个基于 Go + Gin + GORM + MySQL + JWT 的博客系统，包含用户认证、文章管理、评论系统等功能。

## 📋 目录

- [项目简介](#项目简介)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [环境要求](#环境要求)
- [安装步骤](#安装步骤)
- [配置说明](#配置说明)
- [启动项目](#启动项目)
- [API 文档](#api-文档)
- [开发指南](#开发指南)
- [常见问题](#常见问题)

## 🚀 项目简介

这是一个完整的博客系统后端项目，提供以下功能：

- ✅ 用户注册、登录、JWT认证
- ✅ 文章的增删改查
- ✅ 评论系统
- ✅ 用户权限管理
- ✅ RESTful API 设计
- ✅ 数据库自动迁移

## 🛠 技术栈

- **语言**: Go 1.25.4
- **Web框架**: Gin v1.11.0
- **ORM**: GORM v1.31.1
- **数据库**: MySQL 8.0+
- **认证**: JWT (golang-jwt/jwt/v5)
- **其他依赖**:
  - MySQL驱动: gorm.io/driver/mysql
  - 参数验证: go-playground/validator/v10

## 📁 项目结构

```
go-job/
├── main.go                 # 程序入口
├── go.mod                  # Go模块依赖
├── go.sum                  # 依赖版本锁定
├── README.md               # 项目说明文档
├── blog/                   # 博客模块
│   ├── db.go              # 数据库连接和初始化
│   ├── jwt_utils.go       # JWT工具函数
│   ├── middleware.go      # 中间件（认证等）
│   ├── route.go           # 路由定义
│   └── Model/             # 数据模型
│       ├── users.go       # 用户模型
│       ├── posts.go       # 文章模型
│       └── comments.go    # 评论模型
├── task-1/                # 任务1相关代码
├── task-2/                # 任务2相关代码
├── task-3/                # 任务3相关代码
└── utils/                 # 工具函数
    └── fibonacci.go       # 斐波那契数列实现
```

## 💻 环境要求

### 基础环境
- **Go**: 1.25.4 或更高版本
- **MySQL**: 8.0 或更高版本
- **操作系统**: Windows 10/11, macOS, Linux

### 开发工具（推荐）
- **IDE**: VS Code, GoLand, 或其他支持Go的编辑器
- **API测试**: Postman, curl, 或其他HTTP客户端
- **数据库管理**: MySQL Workbench, phpMyAdmin, 或其他MySQL客户端

## 📦 安装步骤

### 1. 克隆项目
```bash
git clone <your-repository-url>
cd go-job
```

### 2. 安装 Go 环境
访问 [Go官网](https://golang.org/dl/) 下载并安装 Go 1.25.4+

验证安装：
```bash
go version
```

### 3. 安装 MySQL
- **Windows**: 下载 [MySQL Installer](https://dev.mysql.com/downloads/installer/)
- **macOS**: 使用 Homebrew: `brew install mysql`
- **Linux**: `sudo apt-get install mysql-server` (Ubuntu/Debian)

### 4. 安装项目依赖
```bash
go mod download
go mod tidy
```

## ⚙️ 配置说明

### 1. 数据库配置

在 `blog/db.go` 文件中修改数据库连接信息：

```go
// 当前配置
mysql.Open("root:123456@(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local")

// 修改为你的配置
mysql.Open("用户名:密码@(主机:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local")
```

### 2. 创建数据库
```sql
CREATE DATABASE gorm CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. JWT 配置

在 `blog/jwt_utils.go` 文件中可以修改JWT配置：

```go
var (
    JWTSecret           = []byte("jss@13&^()")  // 修改为你的密钥
    TokenExpireDuration = time.Hour * 24        // Token过期时间
)
```

## 🚀 启动项目

### 1. 确保 MySQL 服务运行
```bash
# Windows (以管理员身份运行)
net start mysql

# macOS/Linux
sudo systemctl start mysql
# 或
sudo service mysql start
```

### 2. 启动应用
```bash
# 方式1: 直接运行
go run main.go

# 方式2: 编译后运行
go build -o app main.go
./app  # Linux/macOS
app.exe  # Windows
```

### 3. 验证启动
看到以下输出表示启动成功：
```
Database migrated successfully
[GIN-debug] Listening and serving HTTP on :8080
```

访问 http://localhost:8080 验证服务是否正常运行。

## 📚 API 文档

### 认证相关

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/auth/register` | 用户注册 | 否 |
| POST | `/auth/login` | 用户登录 | 否 |
| POST | `/auth/refresh` | 刷新Token | 是 |

### 用户相关

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/profile` | 获取用户信息 | 是 |
| POST | `/api/users` | 创建用户 | 是 |

### 文章相关

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/api/posts` | 创建文章 | 是 |
| GET | `/api/posts/:id` | 获取文章详情 | 是 |
| GET | `/public/posts` | 获取公开文章列表 | 否 |

### 评论相关

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/api/comments` | 创建评论 | 是 |

### 请求示例

#### 用户注册
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456"
  }'
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }'
```

#### 访问受保护资源
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🔧 开发指南

### 添加新的API端点

1. 在 `blog/Model/` 中定义数据模型
2. 在 `blog/route.go` 中添加处理函数
3. 在 `RegisterRoutes` 函数中注册路由

### 添加中间件

在 `blog/middleware.go` 中定义中间件，然后在路由中使用：

```go
// 定义中间件
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
        c.Next()
    }
}

// 使用中间件
r.Use(CustomMiddleware())
```

### 数据库迁移

项目启动时会自动执行数据库迁移。如果需要手动迁移：

```go
db.AutoMigrate(&model.YourNewModel{})
```

## 🐛 常见问题

### 1. 数据库连接失败
**问题**: `Error 1045: Access denied for user 'root'@'localhost'`

**解决方案**:
- 检查MySQL用户名和密码是否正确
- 确保MySQL服务正在运行
- 检查数据库是否存在

### 2. 端口被占用
**问题**: `bind: address already in use`

**解决方案**:
```bash
# 查找占用8080端口的进程
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # macOS/Linux

# 杀死进程或修改端口
```

### 3. Go模块依赖问题
**问题**: `module not found` 或依赖下载失败

**解决方案**:
```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download
go mod tidy

# 如果在中国，设置代理
go env -w GOPROXY=https://goproxy.cn,direct
```

### 4. JWT Token 无效
**问题**: `Invalid token` 错误

**解决方案**:
- 检查Token是否过期
- 确保请求头格式正确: `Authorization: Bearer <token>`
- 检查JWT密钥是否一致

### 5. 数据库表不存在
**问题**: `Table doesn't exist` 错误

**解决方案**:
- 确保数据库迁移成功执行
- 检查数据库连接配置
- 手动运行迁移或重启应用

## 📞 技术支持

如果遇到问题，可以通过以下方式获取帮助：

1. 查看项目日志输出
2. 检查数据库连接状态
3. 验证API请求格式
4. 查看Go和MySQL版本兼容性

## 📄 许可证

本项目仅用于学习和开发目的。

---

**最后更新**: 2025年12月8日
**Go版本**: 1.25.4
**项目版本**: 1.0.0
