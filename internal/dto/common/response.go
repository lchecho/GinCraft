package common

// PaginationResponse 分页响应参数
type PaginationResponse struct {
	Page     int   `json:"page" example:"1"`       // 当前页码
	PageSize int   `json:"page_size" example:"10"` // 每页数量
	Total    int64 `json:"total" example:"100"`    // 总数量
}

// ListResponse 列表响应参数
type ListResponse struct {
	List interface{} `json:"list"` // 列表数据
	PaginationResponse
}

// IDResponse 通用ID响应参数
type IDResponse struct {
	ID interface{} `json:"id" example:"1"` // ID，可以是int、uint、string等
}

// MessageResponse 通用消息响应参数
type MessageResponse struct {
	Message string `json:"message" example:"操作成功"` // 响应消息
}
