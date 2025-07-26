#!/bin/bash

echo "🚀 启动 Chat MVP 应用..."

# 检查是否存在 .env 文件
if [ ! -f .env ]; then
    echo "⚠️  未找到 .env 文件，请创建 .env 文件并设置 OPENAI_API_KEY"
    echo "📝 可以复制 .env.example 文件并修改其中的配置"
    echo ""
    echo "示例命令:"
    echo "cp .env.example .env"
    echo "然后编辑 .env 文件，设置您的 OpenAI API Key"
    echo ""
    exit 1
fi

# 检查 Go 版本
if ! command -v go &> /dev/null; then
    echo "❌ 未找到 Go，请先安装 Go 1.22.5 或更高版本"
    exit 1
fi

# 检查 OPENAI_API_KEY 是否设置
source .env
if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ OPENAI_API_KEY 未设置，请在 .env 文件中设置您的 OpenAI API Key"
    exit 1
fi

echo "✅ 环境检查通过"
echo "🏢  正在构建应用..."

# 构建应用
go build -o chat-mvp .

if [ $? -eq 0 ]; then
    echo "✅ 构建成功"
    echo "🌐 启动服务器..."
    echo "📱 浏览器访问: http://localhost:${PORT:-8080}"
    echo ""
    echo "按 Ctrl+C 停止服务器"
    echo ""
    ./chat-mvp
else
    echo "❌ 构建失败"
    exit 1
fi