# CyBlog 个人博客系统后端模块设计方案

## 一、整体架构
基于现有项目 DDD（领域驱动设计）架构风格，采用分层架构：
```
API层(internal/route) -> 应用层 -> 领域层(internal/domain) -> 基础设施层(infra)
```

## 二、核心模块设计

### 1. 用户认证与权限模块
**功能：
- 用户注册/登录（支持账号密码、OAuth第三方登录）
- JWT Token 认证与刷新
- 用户信息管理（个人资料、头像、密码修改）
- RBAC权限控制（管理员、普通用户、游客权限）

**实现位置：**
- 领域模型：internal/domain/user/
- 仓储接口：pkg/repo/user_repo.go
- 接口：internal/route/v1/user/
- 中间件：internal/route/middleware/auth.go

**技术方案：
- 使用JWT实现无状态认证
- 密码加密使用bcrypt
- Redis存储Token黑名单实现登出功能

---

### 2. 文章管理模块
**功能：
- 文章CRUD（创建、编辑、删除、查询）
- 文章分类管理
- 文章标签管理
- 文章状态管理（草稿、发布、归档）
- Markdown内容解析与HTML渲染
- 文章搜索功能
- 阅读量统计
- 置顶/推荐文章设置

**实现位置：**
- 领域模型：internal/domain/article/
- 仓储接口：pkg/repo/article_repo.go
- 接口：internal/route/v1/article/

**技术方案：
- MongoDB存储文章内容与元数据
- Elasticsearch/Redis实现全文搜索
- 异步任务处理文章内容解析、图片上传、SEO优化

---

### 3. 评论系统模块
**功能：
- 评论提交与回复
- 评论审核机制
- 评论点赞
- 评论删除与举报
- 评论分页查询

**实现位置：**
- 领域模型：internal/domain/comment/
- 仓储接口：pkg/repo/comment_repo.go
- 接口：internal/route/v1/comment/

**技术方案：
- MySQL/MongoDB存储评论数据
- 敏感词过滤功能
- NATS消息队列实现评论通知推送

---

### 4. 文件与资源管理模块
**功能：
- 图片/文件上传
- 资源分类管理
- 静态资源CDN加速
- 图片自动压缩与格式转换

**实现位置：**
- 领域模型：internal/domain/asset/
- 仓储接口：pkg/repo/asset_repo.go
- 接口：internal/route/v1/asset/

**技术方案：
- MinIO作为对象存储服务
- 支持本地存储、阿里云OSS、七牛云等多种存储后端
- 文件MD5校验避免重复上传

---

### 5. 站点配置模块
**功能：
- 博客基础信息配置（站点名称、简介、备案号、联系人）
- SEO配置（Meta标签、站点地图、robots.txt）
- 第三方服务配置（评论系统、统计代码、CDN配置）
- 主题切换配置
- 导航菜单管理

**实现位置：**
- 领域模型：internal/domain/site/
- 仓储接口：pkg/repo/site_repo.go
- 接口：internal/route/v1/site/

**技术方案：
- MySQL存储配置信息
- Redis缓存配置信息提升访问速度

---

### 6. 统计与分析模块
**功能：
- 网站访问量统计
- 文章阅读量统计
- 用户行为分析
- 访问来源统计
- 数据看板接口

**实现位置：**
- 领域模型：internal/domain/stat/
- 仓储接口：pkg/repo/stat_repo.go
- 接口：internal/route/v1/stat/

**技术方案：
- Redis做UV、PV实时统计
- 异步任务统计数据落库MySQL
- 支持按时间维度多维度统计

---

### 7. 通知推送模块
**功能：
- 新评论通知
- 系统消息通知
- 邮件推送
- 短信/微信通知（可选）

**实现位置：**
- 领域模型：internal/domain/notify/
- 服务：pkg/task/notify_task.go
- 接口：internal/route/v1/notify/

**技术方案：
- NATS消息队列实现异步通知解耦
- 支持多种通知渠道可配置

---

## 三、基础支持模块

### 1. 通用中间件
- 请求日志中间件
- 跨域处理中间件
- 限流中间件
- 权限校验中间件
- 参数校验中间件
- 错误统一处理中间件

### 2. 工具类
- 响应结果统一封装
- 分页工具
- 日期时间处理
- 加密解密工具
- Markdown解析工具
- 敏感词过滤工具

## 四、数据库设计概要

### MySQL 存储表结构：
- 用户表（user）
- 分类表（category）
- 标签表（tag）
- 评论表（comment）
- 站点配置表（site_config）
- 统计数据表（stat）

### MongoDB 存储：
- 文章内容表（article）
- 操作日志表（operation_log）

### Redis 缓存：
- Token黑名单
- 站点配置缓存
- 热门文章缓存
- 统计数据缓存
- 接口限流计数器

## 五、接口设计规范
1. 采用RESTful API设计风格
2. 接口版本控制：`/api/v1/xxx`
3. 统一响应格式：
```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "request_id": "xxx"
}
```
4. 错误码统一规范：
   - 0: 成功
   - 1xxxx: 系统错误
   - 2xxxx: 参数错误
   - 3xxxx: 业务错误
   - 4xxxx: 认证/权限错误

## 六、后续实现步骤建议
1. 先实现用户认证与权限模块
2. 实现文章管理核心模块
3. 实现评论系统
4. 实现资源管理模块
5. 实现站点配置与其他辅助模块
6. 性能优化与安全加固