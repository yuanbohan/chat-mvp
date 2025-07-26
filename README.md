# Chat MVP - OpenAI 聊天应用

基于 Gin 框架的简单聊天应用，支持多用户多轮对话功能。

## 项目简介

一个使用 Go + Gin 框架构建的 RESTful API 服务，集成 OpenAI API 实现智能对话功能。项目包含后端 API 服务和前端聊天界面。

## 功能特性

- ✅ RESTful API 接口设计
- ✅ 集成 OpenAI GPT API
- ✅ 多用户会话管理
- ✅ 消息历史记录
- ✅ 响应式聊天界面
- ✅ 多轮对话支持

## 技术栈

### 后端

- **框架**: [Gin](https://gin-gonic.com/) - Go 语言高性能 Web 框架
- **语言**: Go 1.22.5
- **API**: OpenAI GPT API

### 前端

- **技术**: HTML + CSS + JavaScript
- **布局**: 响应式设计
- **交互**: 原生 JavaScript + Fetch API

## 快速开始

### 1. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，设置您的 OpenAI API Key
```

### 2. 启动应用

```bash
# 使用启动脚本（推荐）
./start.sh

# 或者手动启动
go run main.go
```

### 3. 访问应用

打开浏览器访问: `http://localhost:8080`

## 项目结构

```text
chat-mvp/
├── main.go              # 主程序入口
├── handlers/            # API 处理器
│   └── chat.go         # 聊天相关处理逻辑
├── models/             # 数据模型
│   └── message.go      # 消息结构定义
├── services/           # 业务逻辑层
│   └── openai.go       # OpenAI API 调用
├── static/             # 静态文件
│   ├── index.html      # 聊天页面
│   ├── style.css       # 样式文件
│   └── script.js       # 前端逻辑
├── config/             # 配置文件
│   └── config.go       # 应用配置
├── go.mod              # Go 模块文件
├── go.sum              # 依赖版本锁定
└── README.md           # 项目说明
```

## 环境要求

- Go 1.22.5 或更高版本
- OpenAI API Key

## 许可证

本项目采用 MIT 许可证。