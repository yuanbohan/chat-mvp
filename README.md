# Chat MVP - OpenAI 聊天应用

![Go Version](https://img.shields.io/badge/Go-1.22.5+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![OpenAI](https://img.shields.io/badge/OpenAI-GPT--4o-412991?logo=openai&logoColor=white)

基于 Gin 框架构建的智能聊天应用，集成 OpenAI GPT-4o 模型，支持多用户多轮对话功能。

## 🚀 项目简介

这是一个现代化的聊天应用，采用前后端分离架构：

- **后端**：使用 Go + Gin 构建高性能 RESTful API 服务
- **前端**：响应式 Web 界面，支持实时对话交互
- **AI 能力**：集成 [OpenAI Go SDK](https://github.com/openai/openai-go)，使用 GPT-4o 模型提供智能对话

## ✨ 功能特性

- 🏗️ **RESTful API 架构** - 标准化接口设计，易于扩展
- 🤖 **OpenAI GPT-4o 集成** - 先进的 AI 对话能力
- 👥 **多用户支持** - 独立的会话管理机制
- 💬 **多轮对话** - 保持上下文的连续对话
- 📝 **消息历史** - 完整的对话记录存储
- 📱 **响应式界面** - 适配各种设备屏幕
- ⚡ **高性能** - Gin 框架提供卓越性能

## 🛠️ 技术栈

### 后端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| [Go](https://golang.org/) | 1.22.5+ | 高性能编程语言 |
| [Gin](https://gin-gonic.com/) | Latest | 轻量级 Web 框架 |
| [OpenAI Go SDK](https://github.com/openai/openai-go) | Latest | OpenAI 官方 SDK |

### 前端技术

| 技术 | 说明 |
|------|------|
| HTML5 | 语义化标记 |
| CSS3 | 响应式样式设计 |
| JavaScript (ES6+) | 现代前端交互 |
| Fetch API | 异步网络请求 |

## 🚀 快速开始

### 环境要求

- Go 1.22.5 或更高版本
- OpenAI API Key
- 网络连接（访问 OpenAI API）

### 安装步骤

1. **克隆仓库**

```bash
git clone https://github.com/yuanbohan/chat-mvp.git
cd chat-mvp
```

2. **安装依赖**

```bash
go mod tidy
```

3. **配置环境变量**

```bash
# 方式一：使用环境变量
export OPENAI_API_KEY="your-openai-api-key-here"

# 方式二：创建 .env 文件（推荐）
echo "OPENAI_API_KEY=your-openai-api-key-here" > .env
```

4. **启动应用**

```bash
go run main.go
```

5. **访问应用**

打开浏览器访问：`http://127.0.0.1:8080`

## 📁 项目结构

```text
chat-mvp/
├── main.go                 # 🎯 应用程序入口
├── handlers/               # 🎮 API 路由处理器
│   └── chat.go            #    聊天相关接口逻辑
├── models/                # 📊 数据模型定义
│   └── message.go         #    消息结构体
├── services/              # 🔧 业务逻辑服务层
│   └── openai.go          #    OpenAI API 调用封装
├── static/                # 🌐 前端静态资源
│   ├── index.html         #    聊天界面页面
│   ├── style.css          #    界面样式文件
│   └── script.js          #    前端交互逻辑
├── config/                # ⚙️ 应用配置管理
│   └── config.go          #    配置文件处理
├── Dockerfile             # 🐳 Docker 容器化配置
├── go.mod                 # 📦 Go 模块依赖
├── go.sum                 # 🔒 依赖版本锁定
└── README.md              # 📖 项目文档
```

## 🔧 API 接口

### 聊天接口

```http
POST /api/chat
Content-Type: application/json

{
  "message": "你好，请介绍一下你自己",
  "user_id": "user123"
}
```

**响应示例：**

```json
{
  "success": true,
  "message": "你好！我是基于 GPT-4o 的 AI 助手...",
  "timestamp": "2025-07-26T10:30:00Z"
}
```

## 🧪 开发指南

### 本地开发

```bash
go run main.go
```

### 构建部署

```bash
# 构建二进制文件
go build -o chat-mvp main.go

# 运行
./chat-mvp
```

### Docker 部署

使用项目根目录的 `Dockerfile` 进行容器化部署：

```bash
# 构建镜像
docker build -t chat-mvp .

# 运行容器
docker run -d \
  --name chat-mvp \
  -p 8080:8080 \
  -e OPENAI_API_KEY="your-api-key-here" \
  chat-mvp
```

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 更新日志

### v1.0.0 (2025-07-26)

- ✨ 初始版本发布
- 🎯 支持基础聊天功能
- 🤖 集成 OpenAI GPT-4o 模型
- 📱 响应式前端界面

## ❓ 常见问题

**Q: 如何获取 OpenAI API Key？**
A: 访问 [OpenAI Platform](https://platform.openai.com/) 注册账号并创建 API Key。

**Q: 支持哪些 OpenAI 模型？**
A: 当前版本使用 GPT-4o 模型，可在 `services/openai.go` 中修改模型配置。

**Q: 如何修改服务端口？**
A: 在 `main.go` 中修改端口配置或使用环境变量 `PORT`。

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源协议。

## 🔗 相关链接

- [OpenAI API 文档](https://platform.openai.com/docs)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [Go 官方文档](https://golang.org/doc/)

---

**如果这个项目对你有帮助，请给一个 ⭐ Star！**

Made with ❤️ by [yuanbohan](https://github.com/yuanbohan)
