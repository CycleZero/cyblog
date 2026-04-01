# 评论系统实现方案

## 1. 模块概述

评论模块是文章系统的重要补充，支持多级嵌套评论、点赞等功能，遵循项目现有的分层架构设计。

### 1.1 模块职责
| 功能 | 描述 |
|------|------|
| **评论发布** | 用户对文章发表评论，支持一级评论和回复评论 |
| **评论查询** | 获取文章评论列表，支持分页、树形结构展示 |
| **评论管理** | 编辑、删除评论 |
| **评论交互** | 评论点赞、取消点赞 |
| **统计功能** | 文章评论数统计 |

### 1.2 分层架构
评论模块遵循项目统一的分层架构：
- **Model层**：定义Comment数据结构，与数据库表对应
- **Repo层**：评论数据库操作封装
- **Biz层**：核心业务逻辑层，处理评论发布、权限控制等
- **Service层**：HTTP请求处理层，参数校验、调用Biz方法
- **Route层**：HTTP路由定义

## 2. 领域模型设计

### 2.1 核心实体
```go
// Comment 评论模型
type Comment struct {
    gorm.Model
    ArticleID  uint      `gorm:"not null;index;comment:文章ID"`
    UserID     uint      `gorm:"not null;index;comment:评论者ID"`
    ParentID   uint      `gorm:"default:0;index;comment:父评论ID（0表示一级评论）"`
    Content    string    `gorm:"type:text;not null;comment:评论内容"`
    Likes      int       `gorm:"default:0;comment:点赞数"`
    IP         string    `gorm:"size:45;comment:评论者IP"`
    UserAgent  string    `gorm:"size:500;comment:评论者UserAgent"`
    
    // 关联
    Article   *Article   `gorm:"foreignKey:ArticleID"`
    User      *User      `gorm:"foreignKey:UserID"`
    Parent    *Comment   `gorm:"foreignKey:ParentID"`
    Replies   []*Comment `gorm:"foreignKey:ParentID;order:created_at asc"`
}
```

### 2.2 DTO定义
```go
// CreateRequest 创建评论请求
type CreateRequest struct {
    ArticleID uint   `json:"article_id" binding:"required"`
    ParentID  uint   `json:"parent_id"` // 0表示一级评论
    Content   string `json:"content" binding:"required,min=1,max=2000"`
}

// UpdateRequest 更新评论请求
type UpdateRequest struct {
    ID      uint   `json:"id" binding:"required"`
    Content string `json:"content" binding:"required,min=1,max=2000"`
}

// ListQuery 评论列表查询参数
type ListQuery struct {
    common.PageParam
    ArticleID uint   `form:"article_id" binding:"required"`
    ParentID  *uint  `form:"parent_id"` // 不传则获取一级评论，传0也获取一级评论
    SortBy    string `form:"sort_by,default=created_at"` // created_at, likes
    SortOrder string `form:"sort_order,default=desc"`    // asc, desc
}

// Getter方法供Repo层动态查询使用
func (q *ListQuery) GetArticleID() uint   { return q.ArticleID }
func (q *ListQuery) GetParentID() *uint    { return q.ParentID }
func (q *ListQuery) GetSortBy() string     { return q.SortBy }
func (q *ListQuery) GetSortOrder() string  { return q.SortOrder }
func (q *ListQuery) GetPage() int           { return q.Page }
func (q *ListQuery) GetPageSize() int       { return q.PageSize }

// AdminListQuery 管理端评论列表查询参数
type AdminListQuery struct {
    common.PageParam
    Keyword    string `form:"keyword"`    // 搜索评论内容
    ArticleID  uint   `form:"article_id"` // 按文章筛选
    UserID     uint   `form:"user_id"`    // 按用户筛选
    SortBy     string `form:"sort_by,default=created_at"`
    SortOrder  string `form:"sort_order,default=desc"`
}

// Response 评论响应
type Response struct {
    ID        uint      `json:"id"`
    ArticleID uint      `json:"article_id"`
    UserID    uint      `json:"user_id"`
    ParentID  uint      `json:"parent_id"`
    Content   string    `json:"content"`
    Likes     int       `json:"likes"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    User *struct {
        ID     uint   `json:"id"`
        Name   string `json:"name"`
        Avatar string `json:"avatar"`
    } `json:"user"`
    
    // 回复评论时包含被回复的用户信息
    ReplyTo *struct {
        ID     uint   `json:"id"`
        Name   string `json:"name"`
    } `json:"reply_to,omitempty"`
    
    // 子评论（分页查询时一级评论不包含回复，需要单独查询）
    Replies []*Response `json:"replies,omitempty"`
}

// ListResponse 评论列表响应
type ListResponse struct {
    common.PageInfo
    List []*Response `json:"list"`
}

// ArticleCommentCountResponse 文章评论数响应
type ArticleCommentCountResponse struct {
    ArticleID    uint `json:"article_id"`
    CommentCount int  `json:"comment_count"`
}
```

## 3. 分层设计

### 3.1 Repo层接口
```go
type CommentRepo interface {
    Create(ctx context.Context, comment *model.Comment) error
    Update(ctx context.Context, comment *model.Comment) error
    Delete(ctx context.Context, id uint) error
    GetByID(ctx context.Context, id uint) (*model.Comment, error)
    List(ctx context.Context, query interface{}) ([]*model.Comment, int64, error)
    
    // 批量获取文章评论数
    GetCommentCountByArticleIDs(ctx context.Context, articleIDs []uint) (map[uint]int, error)
    
    // 获取单篇文章评论数
    GetCommentCountByArticleID(ctx context.Context, articleID uint) (int, error)
    
    // 获取评论的子评论
    GetReplies(ctx context.Context, parentID uint, page, pageSize int) ([]*model.Comment, int64, error)
    
    // 点赞相关
    IncrementLikes(ctx context.Context, id uint) error
    DecrementLikes(ctx context.Context, id uint) error
    
    // 删除文章的所有评论
    DeleteByArticleID(ctx context.Context, articleID uint) error
}
```

### 3.2 Biz层接口
```go
type CommentBiz interface {
    // 创建评论
    Create(ctx context.Context, req *CreateRequest) (*Response, error)
    
    // 更新评论（仅作者或管理员）
    Update(ctx context.Context, req *UpdateRequest) (*Response, error)
    
    // 删除评论（仅作者或管理员）
    Delete(ctx context.Context, id uint) error
    
    // 获取评论详情
    GetByID(ctx context.Context, id uint) (*Response, error)
    
    // 获取文章评论列表
    List(ctx context.Context, query *ListQuery) (*ListResponse, error)
    
    // 获取评论的回复列表
    GetReplies(ctx context.Context, parentID uint, page, pageSize int) (*ListResponse, error)
    
    // 管理端获取评论列表
    AdminList(ctx context.Context, query *AdminListQuery) (*ListResponse, error)
    
    // 评论点赞
    Like(ctx context.Context, commentID uint) error
    
    // 取消点赞
    Unlike(ctx context.Context, commentID uint) error
    
    // 获取文章评论数
    GetCommentCount(ctx context.Context, articleID uint) (*ArticleCommentCountResponse, error)
    
    // 批量获取文章评论数
    BatchGetCommentCount(ctx context.Context, articleIDs []uint) (map[uint]int, error)
}
```

Biz层内获取用户ID和请求信息方式：
```go
meta := common.GetRequestMetadata(ctx)
userID := meta.UserID
ip := c.ClientIP() // Service层获取IP后传递给Biz或通过context传递
userAgent := c.GetHeader("User-Agent")
```

### 3.3 Service层职责
Service层作为HTTP请求入口，职责包括：
1. 参数绑定与校验
2. 从gin上下文获取IP、UserAgent等信息
3. 调用Biz层方法
4. 统一处理响应

#### Service层示例
```go
type CommentService struct {
    common.BaseService
    biz    *CommentBiz
    logger *log.Logger
}

func (s *CommentService) Create(c *gin.Context) {
    var req CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        s.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
        return
    }
    
    // 可以将IP和UserAgent放入context传递给Biz层
    ctx := context.WithValue(c, "ip", c.ClientIP())
    ctx = context.WithValue(ctx, "user_agent", c.GetHeader("User-Agent"))
    
    resp, err := s.biz.Create(ctx, &req)
    if err != nil {
        s.Error(c, err)
        return
    }
    
    s.Success(c, resp)
}
```

## 4. API接口设计

### 4.1 公开接口（用户端）
| 方法 | 路径 | 权限 | 描述 |
|------|------|------|------|
| POST | /api/v1/comments | 登录用户 | 发表评论 |
| PUT | /api/v1/comments/:id | 作者/管理员 | 更新评论 |
| DELETE | /api/v1/comments/:id | 作者/管理员 | 删除评论 |
| GET | /api/v1/comments | 公开 | 获取文章评论列表 |
| GET | /api/v1/comments/:id/replies | 公开 | 获取评论的回复列表 |
| GET | /api/v1/articles/:id/comment-count | 公开 | 获取文章评论数 |
| POST | /api/v1/comments/:id/like | 登录用户 | 点赞评论 |
| DELETE | /api/v1/comments/:id/like | 登录用户 | 取消点赞 |

### 4.2 管理接口
| 方法 | 路径 | 权限 | 描述 |
|------|------|------|------|
| GET | /api/v1/admin/comments | 管理员 | 获取所有评论列表 |
| DELETE | /api/v1/admin/comments/:id | 管理员 | 删除评论 |

## 5. 数据库表结构设计

### 5.1 comments 评论表
| 字段名 | 类型 | 允许空 | 默认值 | 注释 |
|--------|------|--------|--------|------|
| id | uint | 否 | auto_increment | 主键 |
| article_id | uint | 否 | | 文章ID（索引） |
| user_id | uint | 否 | | 评论者ID（索引） |
| parent_id | uint | 否 | 0 | 父评论ID，0表示一级评论（索引） |
| content | text | 否 | | 评论内容 |
| likes | int | 否 | 0 | 点赞数 |
| ip | varchar(45) | 是 | | 评论者IP地址（支持IPv6） |
| user_agent | varchar(500) | 是 | | 评论者UserAgent |
| created_at | datetime | 否 | | 创建时间 |
| updated_at | datetime | 否 | | 更新时间 |
| deleted_at | datetime | 是 | | 删除时间 |

**索引设计：**
- PRIMARY KEY (id)
- INDEX idx_article_id (article_id)
- INDEX idx_user_id (user_id)
- INDEX idx_parent_id (parent_id)
- INDEX idx_created_at (created_at)
- 联合索引: INDEX idx_article_parent (article_id, parent_id, created_at)

## 6. 业务规则说明

### 6.1 评论发布规则
1. 必须登录才能发表评论
2. 评论内容长度限制：1-2000字符
3. 一级评论parent_id为0，回复评论parent_id为被回复的评论ID
4. 评论发布后立即可见
5. 需记录评论者的IP和UserAgent

### 6.2 评论查询规则
1. 支持分页查询
2. 默认按创建时间倒序排列，也可按点赞数排序
3. 一级评论和回复分开查询：
   - 第一次查询获取一级评论（parent_id=0）
   - 前端点击"查看更多回复"时再查询该评论下的子评论

### 6.3 评论编辑/删除规则
1. 评论作者可以编辑自己的评论
2. 评论作者可以删除自己的评论
3. 管理员可以编辑/删除任意评论
4. 删除评论时，如果该评论有子评论，子评论如何处理？
   - 方案一：级联删除所有子评论（简单）
   - 方案二：子评论保留，父评论显示"该评论已删除"（推荐）

### 6.4 点赞规则
1. 必须登录才能点赞
2. 点赞前需检查评论是否存在
3. TODO：后续添加点赞记录表，防止重复点赞

## 7. 实现步骤

### 7.1 基础层实现
1. 创建`pkg/model/comment.go`，定义Comment模型
2. 在`pkg/repo/article.go`的AutoMigrate中添加Comment模型迁移，或创建`pkg/repo/comment.go`

### 7.2 Repo层实现
3. 创建`pkg/repo/comment.go`，实现CommentRepo接口

### 7.3 Biz层实现
4. 创建`internal/domain/comment`目录：
   - `dto.go`：定义评论相关请求/响应DTO
   - `biz.go`：实现CommentBiz业务逻辑
   - `service.go`：实现CommentService
   - `provider.go`：依赖注入配置
5. 在`internal/domain/provider.go`中注册Comment模块Provider

### 7.4 Route层实现
6. 创建`internal/route/comment.go`，定义评论HTTP路由
7. 在主路由中注册评论路由

### 7.5 集成与测试
8. 在Article模块中集成评论数展示（可选）
9. 编写单元测试
10. 接口测试

## 8. 后续优化方向

1. **敏感词过滤**：评论发布前进行敏感词检测
2. **点赞记录表**：记录用户点赞行为，防止重复点赞
3. **评论通知**：用户收到回复时发送通知
4. **评论举报**：用户可以举报不当评论
5. **富文本支持**：支持简单的Markdown或HTML标签
6. **评论排序策略**：热门评论、最新评论等多种排序方式
7. **评论缓存**：对热点文章的评论进行缓存