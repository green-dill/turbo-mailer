package schema

// Page is a generic pagination structure for API responses
type Page[T any] struct {
	Total       int64  `json:"total"`        // Total number of records
	CurrentPage int    `json:"current_page"` // Current page number
	TotalPages  int    `json:"total_pages"`  // Total number of pages
	List        []T    `json:"list"`         // List of data items
	Next        string `json:"next"`         // URL for the next page
	Prev        string `json:"prev"`         // URL for the previous page
}

// NewPage creates a new instance of Page
func NewPage[T any](total int64, currentPage, pageSize int, list []T, next, prev string) *Page[T] {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &Page[T]{
		Total:       total,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		List:        list,
		Next:        next,
		Prev:        prev,
	}
}
