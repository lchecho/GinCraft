package common

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1" example:"1"`                    // 页码，默认1
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100" example:"10"` // 每页数量，默认10，最大100
}

// GetPage 获取页码，如果为0则返回默认值1
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页数量，如果为0则返回默认值10
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	return p.PageSize
}
