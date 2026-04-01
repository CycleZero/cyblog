package common

// PageParam 分页请求参数
type PageParam struct {
	Page     int `json:"page" uri:"page" binding:"page"`
	PageSize int `json:"page_size" uri:"page_size" binding:"page_size"`
}

// PageInfo 请求结果分页信息
type PageInfo struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
