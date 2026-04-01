package comment

import (
	commonDto "cyblog/internal/common/dto"
	"time"
)

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
	commonDto.PageParam
	ArticleID uint   `form:"article_id" binding:"required"`
	ParentID  *uint  `form:"parent_id"`                  // 不传则获取一级评论，传0也获取一级评论
	SortBy    string `form:"sort_by,default=created_at"` // created_at, likes
	SortOrder string `form:"sort_order,default=desc"`    // asc, desc
}

// Getter方法供Repo层动态查询使用
func (q *ListQuery) GetArticleID() uint   { return q.ArticleID }
func (q *ListQuery) GetParentID() *uint   { return q.ParentID }
func (q *ListQuery) GetSortBy() string    { return q.SortBy }
func (q *ListQuery) GetSortOrder() string { return q.SortOrder }
func (q *ListQuery) GetPage() int         { return q.Page }
func (q *ListQuery) GetPageSize() int     { return q.PageSize }

// AdminListQuery 管理端评论列表查询参数
type AdminListQuery struct {
	commonDto.PageParam
	Keyword   string `form:"keyword"`    // 搜索评论内容
	ArticleID uint   `form:"article_id"` // 按文章筛选
	UserID    uint   `form:"user_id"`    // 按用户筛选
	SortBy    string `form:"sort_by,default=created_at"`
	SortOrder string `form:"sort_order,default=desc"`
}

// Response 评论响应
type Response struct {
	ID        uint      `json:"id"`
	ArticleID uint      `json:"article_id"`
	UserID    uint      `json:"user_id"`
	ParentID  uint      `json:"parent_id"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	IsLiked   bool      `json:"is_liked"` // 当前用户是否已点赞
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User *commonDto.User `json:"user"`

	// 回复评论时包含被回复的用户信息
	ReplyTo *commonDto.User `json:"reply_to,omitempty"`

	// 子评论（分页查询时一级评论不包含回复，需要单独查询）
	Replies []*Response `json:"replies,omitempty"`
}

// ListResponse 评论列表响应
type ListResponse struct {
	commonDto.PageInfo
	List []*Response `json:"list"`
}

// ArticleCommentCountResponse 文章评论数响应
type ArticleCommentCountResponse struct {
	ArticleID    uint `json:"article_id"`
	CommentCount int  `json:"comment_count"`
}
