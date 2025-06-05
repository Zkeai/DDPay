package model

// CreateSubsiteReq 创建分站请求
type CreateSubsiteReq struct {
	Name           string  `json:"name" binding:"required"`
	Subdomain      string  `json:"subdomain" binding:"required"`
	CommissionRate float64 `json:"commission_rate" binding:"required,gte=0,lte=100"`
}

// UpdateSubsiteReq 更新分站请求
type UpdateSubsiteReq struct {
	ID             int64   `json:"id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Domain         string  `json:"domain"`
	Subdomain      string  `json:"subdomain" binding:"required"`
	Logo           string  `json:"logo"`
	Description    string  `json:"description"`
	Theme          string  `json:"theme"`
	Status         int     `json:"status" binding:"oneof=0 1"`
	CommissionRate float64 `json:"commission_rate" binding:"required,gte=0,lte=100"`
}

// SubsiteJsonConfigReq 保存分站JSON配置请求
type SubsiteJsonConfigReq struct {
	SubsiteID int64                  `json:"subsite_id" binding:"required"`
	Config    map[string]interface{} `json:"config" binding:"required"`
}

// GetSubsiteJsonConfigReq 获取分站JSON配置请求
type GetSubsiteJsonConfigReq struct {
	SubsiteID string `form:"subsite_id" binding:"required"`
} 