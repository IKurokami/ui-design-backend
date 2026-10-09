package domain

// Page is the stable response shape used by management screens. Items are
// bounded by pageSize so a large account cannot make the admin UI fetch the
// whole database just to render one screen.
type Page[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}
