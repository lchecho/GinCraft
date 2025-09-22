.PHONY: build run test clean swagger help

# 默认目标
.DEFAULT_GOAL := help

# 构建项目
build:
	@echo "正在构建项目..."
	go build -o server cmd/server/main.go

# 运行项目
run:
	@echo "正在启动服务器..."
	go run cmd/server/main.go

# 运行测试
test:
	@echo "正在运行测试..."
	go test ./...

# 清理构建文件
clean:
	@echo "正在清理构建文件..."
	rm -f server
	rm -rf docs

# 生成Swagger文档
swagger:
	@echo "正在生成Swagger文档..."
	./scripts/swagger.sh

# 安装依赖
deps:
	@echo "正在安装依赖..."
	go mod tidy
	go mod download

# 格式化代码
fmt:
	@echo "正在格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "正在检查代码..."
	golangci-lint run

# 显示帮助信息
help:
	@echo "可用的命令："
	@echo "  build    - 构建项目"
	@echo "  run      - 运行项目"
	@echo "  test     - 运行测试"
	@echo "  clean    - 清理构建文件"
	@echo "  swagger  - 生成Swagger文档"
	@echo "  deps     - 安装依赖"
	@echo "  fmt      - 格式化代码"
	@echo "  lint     - 代码检查"
	@echo "  help     - 显示帮助信息"
