# GinCraft ⚡

> **Crafted with Gin** - 用 Gin 精心打造的 Web 后端框架

GinCraft 是一个基于 Gin 框架的二次封装 Web 后端通用框架，包含了中间件、GORM、全局日志追踪、通用返回、全局 panic 捕获、全局常量、全局错误码、定时任务以及 Swagger API 文档等功能，做到开箱即用。

🎯 **为什么选择 GinCraft？**
- ⚡ **高性能**：基于 Gin 的高性能 HTTP 框架
- 🛠️ **开箱即用**：内置常用中间件和工具
- 📚 **文档完善**：自动生成 Swagger API 文档
- 🎨 **优雅设计**：结构化的 DTO 管理和优雅路由
- 🔧 **易于扩展**：模块化设计，便于功能扩展

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Gin Version](https://img.shields.io/badge/Gin-1.10+-green.svg)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![GinCraft](https://img.shields.io/badge/GinCraft-Crafted%20with%20Gin-orange.svg)](https://github.com/yourusername/gin-craft)

## 特性

- **配置管理**：使用 Viper 加载和管理配置文件
- **日志系统**：集成 Zap，支持文件和控制台输出，支持日志级别和轮转
- **数据库集成**：使用 GORM 连接 MySQL，支持连接池配置
- **中间件**：
  - CORS 跨域处理
  - 请求日志记录
  - 全局异常恢复
  - 认证中间件
  - 限流中间件
- **全局追踪**：为每个请求生成唯一的 Trace ID，方便日志追踪
- **统一响应**：标准化 API 响应格式
- **错误码管理**：集中管理错误码和错误信息
- **定时任务**：使用 cron 库管理定时任务
- **优雅关闭**：支持服务器优雅关闭
- **API 文档**：集成 Swagger 自动生成 API 文档
- **DTO 管理**：结构化的数据传输对象管理
- **优雅路由**：支持参数自动绑定和验证

## 目录结构

```
├── cmd                     # 命令行入口
│   └── server              # 服务器入口
├── config                  # 配置文件
├── docs                    # Swagger 文档（自动生成）
├── internal                # 内部包
│   ├── app                 # 应用初始化
│   ├── constant            # 全局常量
│   ├── controller          # 控制器
│   ├── dao                 # 数据访问对象
│   ├── dto                 # 数据传输对象
│   │   ├── common          # 通用 DTO
│   │   └── user            # 用户模块 DTO
│   ├── middleware          # 中间件
│   ├── model               # 数据模型
│   ├── pkg                 # 内部工具包
│   │   ├── config          # 配置加载
│   │   ├── cron            # 定时任务
│   │   ├── database        # 数据库连接
│   │   └── router          # 优雅路由
│   ├── router              # 路由
│   └── service             # 业务逻辑
├── pkg                     # 公共包
│   ├── logger              # 日志
│   ├── response            # 通用响应
│   ├── tracer              # 全局追踪
│   └── utils               # 工具函数
├── scripts                 # 脚本
│   └── swagger.sh          # Swagger 文档生成脚本
├── Makefile                # 项目构建脚本
└── README.md               # 项目说明文档
```

## 快速开始

### 环境要求

- Go 1.24 或更高版本
- MySQL 5.7 或更高版本
- Git

### 安装

```bash
# 克隆项目
git clone https://github.com/yourusername/gin-craft.git
cd gin-craft

# 安装依赖
go mod tidy

# 或者使用 Makefile
make deps
```

### 配置

1. 复制配置文件模板：
```bash
cp config/config.yaml.example config/config.yaml
```

2. 修改 `config/config.yaml` 中的数据库配置：
```yaml
mysql:
  host: localhost
  port: 3306
  username: your_username
  password: your_password
  database: goframe
```

### 运行

```bash
# 方式一：直接运行
go run cmd/server/main.go

# 方式二：使用 Makefile
make run

# 方式三：构建后运行
make build
./server
```

### 访问 API 文档

启动服务器后，访问 Swagger API 文档：

```
http://localhost:8080/swagger/index.html
```

## 使用指南

### 项目命令

使用 Makefile 提供的便捷命令：

```bash
make help          # 查看所有可用命令
make build         # 构建项目
make run           # 运行项目
make test          # 运行测试
make swagger       # 生成 Swagger 文档
make clean         # 清理构建文件
make deps          # 安装依赖
make fmt           # 格式化代码
make lint          # 代码检查
```

### 添加新模块

#### 1. 创建 DTO

在 `internal/dto/` 下创建新的模块目录，例如 `product/`：

```bash
mkdir -p internal/dto/product
```

创建请求和响应 DTO：

```go
// internal/dto/product/request.go
package product

import "github.com/liuchen/goframe/internal/dto/common"

// ProductCreateRequest 产品创建请求
type ProductCreateRequest struct {
    Name        string  `json:"name" binding:"required" example:"iPhone 15"`
    Description string  `json:"description" example:"最新款iPhone"`
    Price       float64 `json:"price" binding:"required,min=0" example:"999.99"`
}

// ProductListRequest 产品列表请求
type ProductListRequest struct {
    common.PaginationRequest
    Name string `form:"name" json:"name" example:"iPhone"`
}
```

```go
// internal/dto/product/response.go
package product

import "time"

// ProductResponse 产品响应
type ProductResponse struct {
    ID          uint      `json:"id" example:"1"`
    Name        string    `json:"name" example:"iPhone 15"`
    Description string    `json:"description" example:"最新款iPhone"`
    Price       float64   `json:"price" example:"999.99"`
    CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}
```

#### 2. 创建控制器

在 `internal/controller/` 下创建控制器：

```go
// internal/controller/product.go
package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/liuchen/goframe/internal/dto/product"
    "github.com/liuchen/goframe/internal/dto/common"
    "github.com/liuchen/goframe/internal/service"
)

// ProductController 产品控制器
type ProductController struct {
    *BaseController
}

// NewProductController 创建产品控制器
func NewProductController() *ProductController {
    return &ProductController{
        BaseController: NewBaseController(),
    }
}

// Create 创建产品
// @Summary 创建产品
// @Description 创建新产品
// @Tags 产品管理
// @Accept json
// @Produce json
// @Param request body product.ProductCreateRequest true "产品信息"
// @Success 200 {object} common.MessageResponse "创建成功"
// @Router /api/v1/product [post]
func (pc *ProductController) Create(req *product.ProductCreateRequest) (interface{}, error) {
    // 业务逻辑处理
    return common.MessageResponse{Message: "产品创建成功"}, nil
}
```

#### 3. 创建服务

在 `internal/service/` 下创建服务：

```go
// internal/service/product.go
package service

// ProductService 产品服务
type ProductService struct{}

// Create 创建产品
func (ps *ProductService) Create(name, description string, price float64) error {
    // 业务逻辑实现
    return nil
}
```

#### 4. 添加路由

在 `internal/router/router.go` 中添加路由：

```go
// 创建产品控制器实例
productController := controller.NewProductController()

// 产品路由
productGroup := v1.Group("/product")
{
    productGroup.POST("", elegantRouter.WithRequestHandler(productController.Create))
}
```

### 使用优雅路由

框架提供了优雅路由功能，支持参数自动绑定和验证：

```go
// 自动绑定请求参数到结构体
r.POST("/api/v1/user/register", elegantRouter.WithRequestHandler(userController.Register))
```

### 使用中间件

```go
// 单个中间件
authUser := v1.Group("/user", middleware.AuthMiddleware())

// 多个中间件
admin := v1.Group("/admin",
    middleware.AuthMiddleware(),
    middleware.AdminAuthMiddleware())
```

### 使用日志

```go
// 使用 Zap 日志
logger.Info("This is an info log", zap.String("key", "value"))
logger.Error("This is an error log", zap.Error(err))




```

### 使用响应

```go
// 成功响应
response.Success(c, gin.H{
    "message": "操作成功",
    "data": data,
})

// 失败响应
response.Fail(c, constant.PARAM_ERROR, nil)

// 自定义消息的失败响应
response.FailWithMsg(c, constant.SYSTEM_ERROR, "系统发生错误", nil)
```

### 使用 Swagger 文档

#### 生成文档

```bash
# 使用脚本生成
./scripts/swagger.sh

# 或使用 Makefile
make swagger
```

#### 添加 API 注释

为控制器方法添加 Swagger 注释：

```go
// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口，创建新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body user.UserRegisterRequest true "注册信息"
// @Success 200 {object} common.MessageResponse "注册成功"
// @Failure 400 {object} common.MessageResponse "请求参数错误"
// @Router /api/v1/user/register [post]
func (uc *UserController) Register(req *user.UserRegisterRequest) (interface{}, error) {
    // 业务逻辑
}
```

#### 为 DTO 添加示例

```go
type UserRegisterRequest struct {
    Username string `json:"username" binding:"required" example:"john_doe"`
    Password string `json:"password" binding:"required" example:"123456"`
    Email    string `json:"email" binding:"required,email" example:"john@example.com"`
}
```

## API 接口

### 用户管理接口

| 接口 | 方法 | 描述 | 认证 |
|------|------|------|------|
| `/api/v1/user/register` | POST | 用户注册 | 否 |
| `/api/v1/user/login` | POST | 用户登录 | 否 |
| `/api/v1/user/info` | GET | 获取用户信息 | 是 |

### 认证方式

需要认证的接口使用 Bearer Token 认证：

```
Authorization: Bearer {your_token}
```

在 Swagger UI 中：
1. 点击右上角的 "Authorize" 按钮
2. 输入你的 token（不需要 "Bearer " 前缀）
3. 点击 "Authorize" 确认

## 配置说明

配置文件位于 `config/config.yaml`，包含以下配置项：

```yaml
app:
  name: goframe
  mode: debug  # debug, release, test
  port: 8080
  read_timeout: 60
  write_timeout: 60

log:
  level: debug  # debug, info, warn, error, fatal, panic
  filename: logs/app.log
  max_size: 100  # MB
  max_age: 30    # days
  max_backups: 10
  compress: false

mysql:
  host: localhost
  port: 3306
  username: root
  password: password
  database: goframe
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600  # seconds

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 10

trace:
  enabled: true
  type: jaeger
  jaeger:
    service_name: goframe
    collector_endpoint: http://localhost:14268/api/traces
    sampler_type: const
    sampler_param: 1
```

## 开发指南

### 代码规范

1. **命名规范**：
   - 包名使用小写字母
   - 结构体名使用大驼峰命名
   - 方法名使用小驼峰命名
   - 常量使用大写字母和下划线

2. **注释规范**：
   - 所有导出的函数和结构体必须有注释
   - 使用中文注释
   - API 接口必须添加 Swagger 注释

3. **错误处理**：
   - 使用统一的错误码管理
   - 错误信息要清晰明确
   - 记录详细的错误日志

### 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/controller

# 运行测试并显示覆盖率
go test -cover ./...
```

### 代码检查

```bash
# 格式化代码
make fmt

# 代码检查
make lint
```

## 扩展建议

1. **认证与授权**：完善 JWT 认证和 RBAC 权限控制
2. **缓存集成**：集成 Redis 缓存
3. **消息队列**：集成 Kafka 或 RabbitMQ
4. **单元测试**：为各层添加单元测试
5. **CI/CD**：添加 CI/CD 配置文件
6. **监控告警**：集成 Prometheus 和 Grafana
7. **容器化**：添加 Docker 支持
8. **微服务**：支持微服务架构

## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

- 项目主页：[https://github.com/lchecho/gin-craft](https://github.com/lchecho/gin-craft)
- 问题反馈：[https://github.com/lchecho/gin-craft/issues](https://github.com/lchecho/gin-craft/issues)
- 邮箱：gin-craft@cnyy.de

---

如果这个项目对你有帮助，请给它一个 ⭐️