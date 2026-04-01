package dto

// PageParam 分页请求参数
// swagger:model PageParam
type PageParam struct {
	// 页码
	Page int `json:"page" uri:"page" binding:"page"`

	// 每页数量
	PageSize int `json:"page_size" uri:"page_size" binding:"page_size"`
}

// PageInfo 请求结果分页信息
// swagger:model PageInfo
type PageInfo struct {
	// 总数
	Total int `json:"total"`

	// 当前页码
	Page int `json:"page"`

	// 每页数量
	PageSize int `json:"page_size"`
}

type User struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
