# Cyblog 🚀

一个现代化的全栈博客系统，采用前后端分离架构，支持 Markdown 写作、完整的评论互动等功能。

[English](./README_EN.md) | [中文](./README.md)

## ✨ 特性

### 核心功能
- 📝 **Markdown 文章编辑** - 支持富文本编辑，代码高亮，数学公式
- 💬 **评论系统** - 支持多级回复、点赞、@提及
- 🔍 **全文搜索** - 快速检索文章内容
- 🏷️ **分类与标签** - 灵活的文章组织方式
- 👤 **用户系统** - 注册、登录、个人资料管理
- 📊 **管理后台** - 完整的文章、评论、用户管理

### 技术亮点
- ⚡ **前后端分离** - Go + Gin 后端 / Vue 3 + TypeScript 前端
- 🎯 **依赖注入** - 使用 Google Wire 管理依赖
- 📦 **消息队列** - NATS 实现异步任务处理
- 💾 **对象存储** - MinIO 兼容 S3 协议
- 🔐 **JWT 认证** - 安全的身份验证机制
- 📚 **RESTful API** - 标准化的接口设计
- 🧪 **完整测试** - 单元测试 + E2E 测试

## 🏗️ 技术栈

### 后端
![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-v1.12-00ADD8?style=flat-square&logo=gin-gonic&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=flat-square&logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-9.x-DC382D?style=flat-square&logo=redis&logoColor=white)
![NATS](https://img.shields.io/badge/NATS-1.50-27AAE1?style=flat-square&logo=nats&logoColor=white)
![MinIO](https://img.shields.io/badge/MinIO-兼容S3-519DD9?style=flat-square&logo=minio&logoColor=white)
![OpenAI](https://img.shields.io/badge/OpenAI-GPT-412991?style=flat-square&logo=openai&logoColor=white)

### 前端
![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?style=flat-square&logo=vue.js&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-6.x-3178C6?style=flat-square&logo=typescript&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-8.0-646CFF?style=flat-square&logo=vite&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind-4.x-38B2AC?style=flat-square&logo=tailwind-css&logoColor=white)
![Element Plus](https://img.shields.io/badge/Element%20Plus-2.13-409EFF?style=flat-square&logo=element&logoColor=white)
![Pinia](https://img.shields.io/badge/Pinia-3.0-F7B731?style=flat-square&logo=pinia&logoColor=white)

## 📁 项目结构

```
cyblog/
├── main.go              # 应用入口
├── app.go               # 应用配置
├── wire.go              # Wire 依赖注入
├── makefile             # 构建脚本
│
├── internal/            # 内部包
│   ├── domain/          # 领域层
│   │   ├── article/     # 文章领域
│   │   ├── category/    # 分类领域
│   │   ├── comment/     # 评论领域
│   │   └── user/        # 用户领域
│   ├── route/           # 路由层
│   │   ├── article.go
│   │   ├── comment.go
│   │   ├── admin.go
│   │   └── middlewire/  # 中间件
│   ├── common/          # 公共组件
│   └── provider.go      # 服务提供者
│
├── pkg/                 # 公共包
│   ├── model/           # 数据模型
│   ├── repo/            # 数据访问层
│   ├── infra/           # 基础设施
│   ├── log/             # 日志组件
│   ├── llm/             # LLM 集成
│   ├── task/            # 异步任务
│   └── util/            # 工具函数
│
├── conf/                # 配置管理
│   └── viper.go         # Viper 配置
│
├── docs/                # 文档
│   ├── swagger.json     # API 文档
│   └── *.md            # 设计文档
│
└── web/                 # 前端项目
    ├── src/
    │   ├── api/         # API 请求
    │   ├── views/       # 页面视图
    │   ├── components/  # 组件
    │   ├── stores/      # 状态管理
    │   └── router/      # 路由配置
    └── package.json
```

## 🚀 快速开始

### 环境要求

- Go 1.26+
- Node.js 20.19+ / 22.12+
- MySQL 8.0+
- Redis 9.x
- NATS Server (可选)
- MinIO 

### 后端启动

```bash
# 克隆项目
git clone https://github.com/CycleZero/cyblog.git
cd cyblog

# 安装依赖
go mod download

# 初始化 Wire 依赖注入
wire generate ./...

# 生成 Swagger 文档
swag init

# 启动服务
go run main.go
```

### 前端启动

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 使用 Makefile

```bash
# 初始化并构建
make rebuild

# 仅构建
make build

# 生成 Wire 代码
make wire

# 生成 Swagger 文档
make swag
```

## ⚙️ 配置

项目配置文件位于根目录 `config.yaml`，参考 `config.yaml.example` 创建配置文件：

```bash
# 复制配置文件
cp config.yaml.example config.yaml

# 编辑配置
vim config.yaml
```

主要配置项：
- `server.http.port` - HTTP 服务端口
- `data.db.*` - MySQL 数据库连接信息
- `data.redis.*` - Redis 连接配置
- `data.minio.*` - MinIO 对象存储配置
- `llm.*` - LLM API 配置（支持 DeepSeek、OpenAI 等）
- `jwt.secret` - JWT 密钥

## 📖 API 文档

启动服务后访问 Swagger 文档：

```
http://localhost:8080/swagger/index.html
```

[//]: # (## 🧪 测试)

[//]: # ()
[//]: # (```bash)

[//]: # (# 前端单元测试)

[//]: # (cd web && npm run test:unit)

[//]: # (```)

## 📝 功能模块

### 用户模块
- ✅ 用户注册 / 登录
- ✅ JWT 认证
- ✅ 个人资料管理

[//]: # (- ✅ 头像上传)

### 文章模块
- ✅ 文章列表 / 详情
- ✅ Markdown 编辑器
- ✅ 文章搜索
- ✅ 文章点赞
- ✅ 代码高亮

### 分类与标签
- ✅ 分类管理
- ✅ 标签管理
- ✅ 文章关联

### 评论模块
- ✅ 评论列表
- ✅ 发布评论
- ✅ 评论回复
- ✅ 评论点赞
- ✅ AI 智能回复

### 管理后台
- ✅ 文章管理
- ✅ 评论管理
- ✅ 用户管理
- ✅ 数据统计

## 🛠️ 开发指南

### 代码规范

```bash
# 前端代码检查
cd web && npm run lint

# 代码格式化
cd web && npm run format
```

### 添加新的 API

1. 在 `pkg/model/` 下定义数据模型
1. 在 `internal/domain/` 下定义 DTO
2. 在 `pkg/repo/` 下实现数据访问层
3. 在 `internal/domain/*/service.go` 实现业务逻辑
4. 在 `internal/route/` 注册路由
5. 更新 Swagger 文档 `make swag`

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源。

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - 高性能 Web 框架
- [Vue](https://github.com/vuejs/core) - 渐进式 JavaScript 框架
- [Element Plus](https://github.com/element-plus/element-plus) - Vue 3 UI 组件库
- [Tailwind CSS](https://github.com/tailwindlabs/tailwindcss) - 实用优先 CSS 框架

## 📬 联系

如果你有任何问题或建议，欢迎：

- 提交 [Issue](https://github.com/CycleZero/cyblog/issues)
- 发送 Pull Request
- 联系开发者

---

<p align="center">
  <strong>如果你觉得这个项目有帮助，请点个 ⭐ Star 支持一下！</strong>
</p>
