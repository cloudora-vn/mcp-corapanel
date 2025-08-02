package types

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PageResult struct {
	Total int64 `json:"total"`
}

type PageRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Name     string `json:"name"`
}
