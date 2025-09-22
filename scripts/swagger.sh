#!/bin/bash

# Swagger文档生成脚本

echo "正在生成Swagger文档..."

# 检查swag是否安装
if ! command -v swag &> /dev/null; then
    echo "swag命令未找到，正在安装..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# 生成swagger文档
swag init -g cmd/server/main.go -o docs

echo "Swagger文档生成完成！"
echo "访问地址: http://localhost:8080/swagger/index.html"
