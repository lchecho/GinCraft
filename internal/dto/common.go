package dto

// Pagination 分页请求参数
type Pagination struct {
	NowPage int   `form:"now_page" json:"now_page" binding:"omitempty,min=1" example:"1"`          // 页码，默认1
	PerPage int   `form:"per_page" json:"per_page" binding:"omitempty,min=1,max=100" example:"10"` // 每页数量，默认10，最大100
	Total   int64 `json:"total"`
}

// GetPage 获取页码，如果为0则返回默认值1
func (p *Pagination) GetPage() int {
	if p.NowPage <= 0 {
		return 1
	}
	return p.NowPage
}

// GetPageSize 获取每页数量，如果为0则返回默认值10
func (p *Pagination) GetPageSize() int {
	if p.PerPage <= 0 {
		return 10
	}
	return p.PerPage
}

// IDResponse 通用ID响应参数
type IDResponse struct {
	ID interface{} `json:"id" example:"1"` // ID，可以是int、uint、string等
}
