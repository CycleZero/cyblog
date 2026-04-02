# Cyblog 🚀

A modern full-stack blogging system with a separated frontend and backend architecture, supporting Markdown writing, and complete comment interactions.

[English](./README_EN.md) | [中文](./README.md)

## ✨ Features

### Core Features
- 📝 **Markdown Article Editor** - Rich text editing, code highlighting, math formulas
- 💬 **Comment System** - Multi-level replies, likes, @mentions
- 🔍 **Full-text Search** - Fast article content retrieval
- 🏷️ **Categories & Tags** - Flexible article organization
- 👤 **User System** - Registration, login, profile management
- 📊 **Admin Dashboard** - Complete article, comment, and user management

### Technical Highlights
- ⚡ **Separated Architecture** - Go + Gin backend / Vue 3 + TypeScript frontend
- 🎯 **Dependency Injection** - Powered by Google Wire
- 📦 **Message Queue** - NATS for async task processing
- 💾 **Object Storage** - MinIO with S3 compatibility
- 🔐 **JWT Authentication** - Secure identity verification
- 📚 **RESTful API** - Standardized interface design
- 🧪 **Complete Testing** - Unit tests + E2E tests

## 🏗️ Tech Stack

### Backend
![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-v1.12-00ADD8?style=flat-square&logo=gin-gonic&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=flat-square&logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-9.x-DC382D?style=flat-square&logo=redis&logoColor=white)
![NATS](https://img.shields.io/badge/NATS-1.50-27AAE1?style=flat-square&logo=nats&logoColor=white)
![MinIO](https://img.shields.io/badge/MinIO-S3%20Compatible-519DD9?style=flat-square&logo=minio&logoColor=white)
![OpenAI](https://img.shields.io/badge/DeepSeek/LLM-412991?style=flat-square&logo=openai&logoColor=white)

### Frontend
![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?style=flat-square&logo=vue.js&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-6.x-3178C6?style=flat-square&logo=typescript&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-8.0-646CFF?style=flat-square&logo=vite&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind-4.x-38B2AC?style=flat-square&logo=tailwind-css&logoColor=white)
![Element Plus](https://img.shields.io/badge/Element%20Plus-2.13-409EFF?style=flat-square&logo=element&logoColor=white)
![Pinia](https://img.shields.io/badge/Pinia-3.0-F7B731?style=flat-square&logo=pinia&logoColor=white)

## 📁 Project Structure

```
cyblog/
├── main.go              # Application entry point
├── app.go               # Application configuration
├── wire.go              # Wire dependency injection
├── makefile             # Build scripts
│
├── internal/            # Internal packages
│   ├── domain/          # Domain layer
│   │   ├── article/     # Article domain
│   │   ├── category/   # Category domain
│   │   ├── comment/    # Comment domain
│   │   └── user/       # User domain
│   ├── route/           # Route layer
│   │   ├── article.go
│   │   ├── comment.go
│   │   ├── admin.go
│   │   └── middlewire/ # Middleware
│   ├── common/         # Common components
│   └── provider.go     # Service provider
│
├── pkg/                 # Public packages
│   ├── model/           # Data models
│   ├── repo/            # Data access layer
│   ├── infra/           # Infrastructure
│   ├── log/             # Logging
│   ├── llm/             # LLM integration
│   ├── task/            # Async tasks
│   └── util/            # Utilities
│
├── conf/                # Configuration
│   └── viper.go         # Viper config
│
├── docs/                # Documentation
│   ├── swagger.json     # API docs
│   └── *.md             # Design docs
│
└── web/                 # Frontend project
    ├── src/
    │   ├── api/         # API requests
    │   ├── views/       # Page views
    │   ├── components/  # Components
    │   ├── stores/      # State management
    │   └── router/      # Router config
    └── package.json
```

## 🚀 Quick Start

### Requirements

- Go 1.26+
- Node.js 20.19+ / 22.12+
- MySQL 8.0+
- Redis 9.x
- NATS Server (optional)
- MinIO (optional)

### Backend Setup

```bash
# Clone the project
git clone https://github.com/CycleZero/cyblog.git
cd cyblog

# Install dependencies
go mod download

# Initialize Wire dependency injection
wire generate ./...

# Generate Swagger docs
swag init

# Start the server
go run main.go
```

### Frontend Setup

```bash
cd web

# Install dependencies
npm install

# Start dev server
npm run dev
```

### Using Makefile

```bash
# Initialize and build
make rebuild

# Build only
make build

# Generate Wire code
make wire

# Generate Swagger docs
make swag
```

## ⚙️ Configuration

The configuration file is located at `config.yaml` in the root directory. Create your config from the example:

```bash
# Copy the config file
cp config.yaml.example config.yaml

# Edit the config
vim config.yaml
```

Key configuration options:
- `server.http.port` - HTTP server port
- `data.db.*` - MySQL database connection
- `data.redis.*` - Redis connection config
- `data.minio.*` - MinIO object storage config
- `llm.*` - LLM API config (supports DeepSeek, OpenAI, etc.)
- `jwt.secret` - JWT secret key

## 📖 API Documentation

After starting the server, access Swagger docs at:

```
http://localhost:8080/swagger/index.html
```

## 📝 Features

### User Module
- ✅ User registration / login
- ✅ JWT authentication
- ✅ Profile management

### Article Module
- ✅ Article list / details
- ✅ Markdown editor
- ✅ Article search
- ✅ Article likes
- ✅ Code highlighting

### Categories & Tags
- ✅ Category management
- ✅ Tag management
- ✅ Article associations

### Comment Module
- ✅ Comment list
- ✅ Post comments
- ✅ Comment replies
- ✅ Comment likes
- ✅ AI-powered replies

### Admin Dashboard
- ✅ Article management
- ✅ Comment management
- ✅ User management
- ✅ Statistics

## 🛠️ Development Guide

### Code Standards

```bash
# Frontend linting
cd web && npm run lint

# Code formatting
cd web && npm run format
```

### Adding New API

1. Define data models in `pkg/model/`
2. Define DTOs in `internal/domain/`
3. Implement data access layer in `pkg/repo/`
4. Implement business logic in `internal/domain/*/service.go`
5. Register routes in `internal/route/`
6. Update Swagger docs: `make swag`

## 📄 License

This project is open source under the [MIT License](LICENSE).

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - High-performance web framework
- [Vue](https://github.com/vuejs/core) - Progressive JavaScript framework
- [Element Plus](https://github.com/element-plus/element-plus) - Vue 3 UI component library
- [Tailwind CSS](https://github.com/tailwindlabs/tailwindcss) - Utility-first CSS framework

## 📬 Contact

For questions or suggestions:

- Submit an [Issue](https://github.com/CycleZero/cyblog/issues)
- Send a Pull Request
- Contact the developer

---

<p align="center">
  <strong>If you find this project helpful, please give it a ⭐ Star!</strong>
</p>
