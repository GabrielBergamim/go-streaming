package video

type Page[T any] struct {
	TotalItems  int64 `json:"totalItems"`
	TotalPages  int64   `json:"totalPages"`
	IsLast bool `json:"isLast"`
    Content []T `json:"content"`
}
