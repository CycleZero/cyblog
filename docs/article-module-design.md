# 文章模块实现方案

## 1. 模块概述
文章相关模块分为三个独立的子模块，每个模块遵循项目现有的分层架构，职责单一，互相解耦：

### 1.1 模块划分
| 模块 | 职责 |
|------|------|
| **Article模块** | 文章核心管理：创建、编辑、查询、状态管理、点赞、浏览统计等 |
| **Category模块** | 分类管理：分类的增删改查、层级管理、排序等 |
| **Tag模块** | 标签管理：标签的增删改查、标签计数、颜色管理等 |

### 1.2 分层架构
每个子模块都有独立的完整分层，各层职责严格划分：
- **Model层**：定义对应的数据结构（PO），与数据库表一一对应
- **Repo层**：对应模块的数据库操作封装，只做CRUD，无业务逻辑
- **Biz层**：核心业务逻辑层，处理所有业务规则、事务、跨模块调用
- **Service层**：HTTP请求处理层，仅做参数校验、上下文解析、调用Biz方法、返回响应
- **Route层**：HTTP路由定义，绑定URL路径与Service方法

## 2. 领域模型设计

### 2.1 核心实体
```go
// Article 文章模型
type Article struct {
    gorm.Model
    Title         string   `gorm:"size:255;not null;comment:文章标题"`
    Slug          string   `gorm:"size:255;unique;comment:文章别名（用于URL）"`
    Content       string   `gorm:"type:text;not null;comment:文章内容（Markdown）"`
    HtmlContent   string   `gorm:"type:text;not null;comment:渲染后的HTML内容"`
    Summary       string   `gorm:"size:500;comment:文章摘要"`
    CoverImage    string   `gorm:"size:255;comment:封面图URL"`
    Status        int      `gorm:"default:1;comment:状态：1草稿 2已发布 3已归档"`
    Views         int      `gorm:"default:0;comment:浏览量"`
    Likes         int      `gorm:"default:0;comment:点赞数"`
    IsTop         bool     `gorm:"default:false;comment:是否置顶"`
    IsOriginal    bool     `gorm:"default:true;comment:是否原创"`
    AuthorID      uint     `gorm:"not null;comment:作者ID"`
    CategoryID    uint     `gorm:"comment:分类ID"`
    
    // 关联
    Author        *User    `gorm:"foreignKey:AuthorID"`
    Category      *Category `gorm:"foreignKey:CategoryID"`
    Tags          []*Tag   `gorm:"many2many:article_tags;"`
}

// Category 分类模型
type Category struct {
    gorm.Model
    Name          string   `gorm:"size:50;not null;unique;comment:分类名称"`
    Slug          string   `gorm:"size:50;unique;comment:分类别名"`
    Description   string   `gorm:"size:200;comment:分类描述"`
    ParentID      uint     `gorm:"default:0;comment:父分类ID"`
    Sort          int      `gorm:"default:0;comment:排序"`
    
    // 关联
    Articles      []*Article `gorm:"foreignKey:CategoryID"`
}

// Tag 标签模型
type Tag struct {
    gorm.Model
    Name          string   `gorm:"size:50;not null;unique;comment:标签名称"`
    Slug          string   `gorm:"size:50;unique;comment:标签别名"`
    Color         string   `gorm:"size:7;default:#165DFF;comment:标签颜色"`
    Count         int      `gorm:"default:0;comment:关联文章数"`
    
    // 关联
    Articles      []*Article `gorm:"many2many:article_tags;"`
}
```

### 2.2 DTO定义
```go
// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
    Title         string   `json:"title" binding:"required,max=255"`
    Content       string   `json:"content" binding:"required"`
    Summary       string   `json:"summary" binding:"max=500"`
    CoverImage    string   `json:"cover_image" binding:"url"`
    Status        int      `json:"status" binding:"oneof=1 2 3"`
    IsTop         bool     `json:"is_top"`
    IsOriginal    bool     `json:"is_original"`
    CategoryID    uint     `json:"category_id"`
    TagIDs        []uint   `json:"tag_ids"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
    ID            uint     `json:"id" binding:"required"`
    Title         string   `json:"title" binding:"max=255"`
    Content       string   `json:"content"`
    Summary       string   `json:"summary" binding:"max=500"`
    CoverImage    string   `json:"cover_image" binding:"omitempty,url"`
    Status        int      `json:"status" binding:"omitempty,oneof=1 2 3"`
    IsTop         *bool    `json:"is_top"`
    IsOriginal    *bool    `json:"is_original"`
    CategoryID    *uint    `json:"category_id"`
    TagIDs        []uint   `json:"tag_ids"`
}

// ArticleListQuery 文章列表查询参数
// 嵌入公共分页参数PageParam
type ArticleListQuery struct {
    common.PageParam
    Keyword       string   `form:"keyword"`
    CategoryID    uint     `form:"category_id"`
    TagID         uint     `form:"tag_id"`
    Status        int      `form:"status"`
    AuthorID      uint     `form:"author_id"`
    IsTop         *bool    `form:"is_top"`
    SortBy        string   `form:"sort_by,default=created_at"` // created_at, views, likes
    SortOrder     string   `form:"sort_order,default=desc"`   // asc, desc
}

// ArticleListResponse 文章列表响应
// 嵌入公共分页信息PageInfo
type ArticleListResponse struct {
    common.PageInfo
    List []*ArticleResponse `json:"list"`
}

// ArticleResponse 文章响应
type ArticleResponse struct {
    ID            uint      `json:"id"`
    Title         string    `json:"title"`
    Slug          string    `json:"slug"`
    Summary       string    `json:"summary"`
    CoverImage    string    `json:"cover_image"`
    Status        int       `json:"status"`
    Views         int       `json:"views"`
    Likes         int       `json:"likes"`
    IsTop         bool      `json:"is_top"`
    IsOriginal    bool      `json:"is_original"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    Author        *UserResponse `json:"author"`
    Category      *CategoryResponse `json:"category"`
    Tags          []*TagResponse `json:"tags"`
}
```

## 3. 分层设计

### 3.1 Repo层接口
#### 3.1.1 ArticleRepo 文章数据操作接口
```go
type ArticleRepo interface {
    Create(ctx context.Context, article *model.Article) error
    Update(ctx context.Context, article *model.Article) error
    Delete(ctx context.Context, id uint) error
    GetByID(ctx context.Context, id uint) (*model.Article, error)
    GetBySlug(ctx context.Context, slug string) (*model.Article, error)
    List(ctx context.Context, query *articleDto.ListQuery) ([]*model.Article, int64, error)
    
    // 关联操作
    UpdateArticleTags(ctx context.Context, articleID uint, tagIDs []uint) error
    GetArticleTags(ctx context.Context, articleID uint) ([]*model.Tag, error)
    
    // 统计操作
    IncrementViews(ctx context.Context, id uint) error
    IncrementLikes(ctx context.Context, id uint) error
    DecrementLikes(ctx context.Context, id uint) error
}
```

#### 3.1.2 CategoryRepo 分类数据操作接口
```go
type CategoryRepo interface {
    Create(ctx context.Context, category *model.Category) error
    Update(ctx context.Context, category *model.Category) error
    Delete(ctx context.Context, id uint) error
    GetByID(ctx context.Context, id uint) (*model.Category, error)
    GetBySlug(ctx context.Context, slug string) (*model.Category, error)
    List(ctx context.Context, page, pageSize int) ([]*model.Category, int64, error)
    GetArticleCount(ctx context.Context, categoryID uint) (int64, error)
}
```

#### 3.1.3 TagRepo 标签数据操作接口
```go
type TagRepo interface {
    Create(ctx context.Context, tag *model.Tag) error
    Update(ctx context.Context, tag *model.Tag) error
    Delete(ctx context.Context, id uint) error
    GetByID(ctx context.Context, id uint) (*model.Tag, error)
    GetBySlug(ctx context.Context, slug string) (*model.Tag, error)
    List(ctx context.Context, page, pageSize int) ([]*model.Tag, int64, error)
    GetByIDs(ctx context.Context, ids []uint) ([]*model.Tag, error)
    IncrementCount(ctx context.Context, id uint) error
    DecrementCount(ctx context.Context, id uint) error
}
```

### 3.2 Biz层接口
Biz层方法接收`context.Context`作为上下文参数，不耦合Gin框架，保证通用性：
> 注意：`*gin.Context`实现了`context.Context`接口，可以直接传入Biz层方法
```go
type ArticleBiz interface {
    Create(ctx context.Context, req *articleDto.CreateRequest) (*articleDto.Response, error)
    Update(ctx context.Context, req *articleDto.UpdateRequest) (*articleDto.Response, error)
    Delete(ctx context.Context, id uint) error
    GetByID(ctx context.Context, id uint, incrementViews bool) (*articleDto.Response, error)
    GetBySlug(ctx context.Context, slug string, incrementViews bool) (*articleDto.Response, error)
    List(ctx context.Context, query *articleDto.ListQuery) (*articleDto.ListResponse, error)
    
    // 交互功能
    Like(ctx context.Context, articleID uint) error
    Unlike(ctx context.Context, articleID uint) error
}
```

Biz层内获取用户ID方式：
```go
meta := common.GetRequestMetadata(ctx)
userID := meta.UserID
```

#### 3.2.2 CategoryBiz 分类业务接口
```go
type CategoryBiz interface {
    Create(ctx context.Context, req *categoryDto.CreateRequest) (*categoryDto.Response, error)
    Update(ctx context.Context, req *categoryDto.UpdateRequest) (*categoryDto.Response, error)
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, page, pageSize int) (*categoryDto.ListResponse, error)
}
```

#### 3.2.3 TagBiz 标签业务接口
```go
type TagBiz interface {
    Create(ctx context.Context, req *tagDto.CreateRequest) (*tagDto.Response, error)
    Update(ctx context.Context, req *tagDto.UpdateRequest) (*tagDto.Response, error)
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, page, pageSize int) (*tagDto.ListResponse, error)
}
```

### 3.3 Service层职责
Service层作为HTTP请求入口，只负责以下逻辑，不包含任何业务处理：
1. 参数绑定与校验（使用gin的ShouldBind系列方法）
2. 从gin上下文获取公共信息（如登录用户ID、请求元数据等）
3. 调用**单个Biz方法**完成业务操作（不调用多个Biz，不做业务逻辑编排）
4. 统一处理Biz返回的结果和错误，封装成JSON格式返回给前端

#### Service层示例（ArticleService）
```go
type ArticleService struct {
    common.BaseService
    biz *ArticleBiz
    logger *log.Logger
}

func NewArticleService(biz *ArticleBiz, logger *log.Logger) *ArticleService {
    return &ArticleService{
        biz: biz,
        logger: logger,
    }
}

// CreateArticle 创建文章接口处理
func (s *ArticleService) Create(c *gin.Context) {
    // 1. 参数绑定与校验
    var req articleDto.CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        s.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
        return
    }

    // 2. 直接将gin.Context传递给Biz层（*gin.Context实现了context.Context接口）
    resp, err := s.biz.Create(c, &req)
    if err != nil {
        s.Error(c, err)
        return
    }

    // 3. 返回成功响应
    s.Success(c, resp)
}
```

### 3.4 各层依赖关系
- Service层 → 依赖对应模块的Biz层（仅依赖自身模块Biz）
- Biz层 → 依赖自身模块的Repo层 + 其他模块的Biz层（跨模块调用通过Biz接口）
- Repo层 → 仅依赖数据库ORM，不依赖其他层

## 4. API接口设计

### 4.1 文章相关接口
| 方法 | 路径 | 权限 | 描述 |
|------|------|------|------|
| POST | /api/v1/articles | 登录用户 | 创建文章 |
| PUT | /api/v1/articles/:id | 作者/管理员 | 更新文章 |
| DELETE | /api/v1/articles/:id | 作者/管理员 | 删除文章 |
| GET | /api/v1/articles | 公开 | 获取文章列表（分页） |
| GET | /api/v1/articles/:id | 公开 | 获取文章详情 |
| GET | /api/v1/articles/slug/:slug | 公开 | 根据别名获取文章 |
| POST | /api/v1/articles/:id/like | 登录用户 | 点赞文章 |
| DELETE | /api/v1/articles/:id/like | 登录用户 | 取消点赞 |

### 4.2 分类相关接口
| 方法 | 路径 | 权限 | 描述 |
|------|------|------|------|
| POST | /api/v1/categories | 管理员 | 创建分类 |
| PUT | /api/v1/categories/:id | 管理员 | 更新分类 |
| DELETE | /api/v1/categories/:id | 管理员 | 删除分类 |
| GET | /api/v1/categories | 公开 | 获取分类列表（分页） |

### 4.3 标签相关接口
| 方法 | 路径 | 权限 | 描述 |
|------|------|------|------|
| POST | /api/v1/tags | 管理员 | 创建标签 |
| PUT | /api/v1/tags/:id | 管理员 | 更新标签 |
| DELETE | /api/v1/tags/:id | 管理员 | 删除标签 |
| GET | /api/v1/tags | 公开 | 获取标签列表（分页） |

### 4.4 分页规范说明
所有列表接口统一使用公共分页结构：
- 请求DTO嵌入`common.PageParam`，包含`page`和`page_size`参数
- 响应DTO嵌入`common.PageInfo`，包含`total`、`page`、`page_size`字段，以及`list`字段存储实际数据列表

## 5. 数据库表结构设计

### 5.1 articles 文章表
| 字段名 | 类型 | 允许空 | 默认值 | 注释 |
|--------|------|--------|--------|------|
| id | uint | 否 | auto_increment | 主键 |
| title | varchar(255) | 否 | | 文章标题 |
| slug | varchar(255) | 否 | | 文章别名 |
| content | text | 否 | | 文章内容（Markdown） |
| html_content | text | 否 | | 渲染后的HTML |
| summary | varchar(500) | 是 | | 摘要 |
| cover_image | varchar(255) | 是 | | 封面图 |
| status | tinyint | 否 | 1 | 状态：1草稿 2已发布 3已归档 |
| views | int | 否 | 0 | 浏览量 |
| likes | int | 否 | 0 | 点赞数 |
| is_top | tinyint(1) | 否 | 0 | 是否置顶 |
| is_original | tinyint(1) | 否 | 1 | 是否原创 |
| author_id | uint | 否 | | 作者ID |
| category_id | uint | 是 | | 分类ID |
| created_at | datetime | 否 | | 创建时间 |
| updated_at | datetime | 否 | | 更新时间 |
| deleted_at | datetime | 是 | | 删除时间 |

### 5.2 categories 分类表
| 字段名 | 类型 | 允许空 | 默认值 | 注释 |
|--------|------|--------|--------|------|
| id | uint | 否 | auto_increment | 主键 |
| name | varchar(50) | 否 | | 分类名称 |
| slug | varchar(50) | 否 | | 分类别名 |
| description | varchar(200) | 是 | | 分类描述 |
| parent_id | uint | 否 | 0 | 父分类ID |
| sort | int | 否 | 0 | 排序 |
| created_at | datetime | 否 | | 创建时间 |
| updated_at | datetime | 否 | | 更新时间 |
| deleted_at | datetime | 是 | | 删除时间 |

### 5.3 tags 标签表
| 字段名 | 类型 | 允许空 | 默认值 | 注释 |
|--------|------|--------|--------|------|
| id | uint | 否 | auto_increment | 主键 |
| name | varchar(50) | 否 | | 标签名称 |
| slug | varchar(50) | 否 | | 标签别名 |
| color | varchar(7) | 否 | #165DFF | 标签颜色 |
| count | int | 否 | 0 | 关联文章数 |
| created_at | datetime | 否 | | 创建时间 |
| updated_at | datetime | 否 | | 更新时间 |
| deleted_at | datetime | 是 | | 删除时间 |

### 5.4 article_tags 文章标签关联表
| 字段名 | 类型 | 允许空 | 默认值 | 注释 |
|--------|------|--------|--------|------|
| article_id | uint | 否 | | 文章ID |
| tag_id | uint | 否 | | 标签ID |
| created_at | datetime | 否 | | 创建时间 |
| 主键 | (article_id, tag_id) | | | 联合主键 |

## 6. 实现步骤

### 6.1 基础层实现
1. 创建`pkg/model/article.go`, `pkg/model/category.go`, `pkg/model/tag.go`，定义Article、Category、Tag模型
2. 使用AutoMigrate迁移数据库

### 6.2 Category模块实现
3. 创建`pkg/repo/category.go`，实现CategoryRepo
4. 创建`internal/domain/category`目录：
   - `dto.go`：定义分类相关请求/响应DTO
   - `biz.go`：实现CategoryBiz业务逻辑
   - `service.go`：实现CategoryService
   - `provider.go`：依赖注入配置
5. 在`internal/domain/provider.go`中注册Category模块Provider
6. 创建`internal/route/category.go`，定义分类HTTP路由

### 6.3 Tag模块实现
7. 创建`pkg/repo/tag.go`，实现TagRepo
8. 创建`internal/domain/tag`目录：
   - `dto.go`：定义标签相关请求/响应DTO
   - `biz.go`：实现TagBiz业务逻辑
   - `service.go`：实现TagService
   - `provider.go`：依赖注入配置
9. 在`internal/domain/provider.go`中注册Tag模块Provider
10. 创建`internal/route/tag.go`，定义标签HTTP路由

### 6.4 Article模块实现
11. 创建`pkg/repo/article.go`，实现ArticleRepo
12. 创建`internal/domain/article`目录：
    - `dto.go`：定义文章相关请求/响应DTO
    - `biz.go`：实现ArticleBiz业务逻辑（依赖CategoryBiz、TagBiz）
    - `service.go`：实现ArticleService
    - `provider.go`：依赖注入配置
13. 在`internal/domain/provider.go`中注册Article模块Provider
14. 创建`internal/route/article.go`，定义文章HTTP路由

### 6.5 收尾工作
15. 编写各模块单元测试
16. 接口测试
