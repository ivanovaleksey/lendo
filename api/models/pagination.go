package models

// swagger:parameters paginationParams
type PaginationParams struct {
	// Pagination limit
	// in: query
	Limit int `json:"limit"`
	// Pagination offset
	// in: query
	Offset int `json:"offset"`
}

func (params PaginationParams) GetLimit() int {
	const defaultLimit = 10

	if params.Limit > 0 {
		return params.Limit
	}
	return defaultLimit
}
