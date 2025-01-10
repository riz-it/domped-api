package dto

// Response
type ApiResponse[T any] struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	Data    *T                  `json:"data"`
	Paging  *PageMetadata       `json:"paging,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

func (e *ApiResponse[T]) Error() string {
	return e.Message
}

type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetadata PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}
