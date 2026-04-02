package dto

// AdminArticleQuery 管理端文章查询参数
type AdminArticleQuery struct {
	PageParam
	Keyword    string `form:"keyword" query:"keyword"`
	CategoryID uint   `form:"category_id" query:"category_id"`
	TagID      uint   `form:"tag_id" query:"tag_id"`
	Status     int    `form:"status" query:"status"`
	AuthorID   uint   `form:"author_id" query:"author_id"`
	IsTop      *bool  `form:"is_top" query:"is_top"`
	SortBy     string `form:"sort_by,default=created_at" query:"sort_by"`
	SortOrder  string `form:"sort_order,default=desc" query:"sort_order"`
}

func (q *AdminArticleQuery) GetKeyword() string   { return q.Keyword }
func (q *AdminArticleQuery) GetCategoryID() uint  { return q.CategoryID }
func (q *AdminArticleQuery) GetTagID() uint       { return q.TagID }
func (q *AdminArticleQuery) GetStatus() int       { return q.Status }
func (q *AdminArticleQuery) GetAuthorID() uint    { return q.AuthorID }
func (q *AdminArticleQuery) GetIsTop() *bool      { return q.IsTop }
func (q *AdminArticleQuery) GetSortBy() string    { return q.SortBy }
func (q *AdminArticleQuery) GetSortOrder() string { return q.SortOrder }
func (q *AdminArticleQuery) GetPage() int         { return q.Page }
func (q *AdminArticleQuery) GetPageSize() int     { return q.PageSize }

// AdminUserQuery 管理端用户查询参数
type AdminUserQuery struct {
	PageParam
	Keyword   string `form:"keyword" query:"keyword"`
	Role      string `form:"role" query:"role"`
	Status    int    `form:"status" query:"status"`
	SortBy    string `form:"sort_by,default=created_at" query:"sort_by"`
	SortOrder string `form:"sort_order,default=desc" query:"sort_order"`
}

func (q *AdminUserQuery) GetKeyword() string   { return q.Keyword }
func (q *AdminUserQuery) GetRole() string      { return q.Role }
func (q *AdminUserQuery) GetStatus() int       { return q.Status }
func (q *AdminUserQuery) GetSortBy() string    { return q.SortBy }
func (q *AdminUserQuery) GetSortOrder() string { return q.SortOrder }
func (q *AdminUserQuery) GetPage() int         { return q.Page }
func (q *AdminUserQuery) GetPageSize() int     { return q.PageSize }
