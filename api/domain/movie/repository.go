package movie

import "github.com/google/uuid"

type MovieRepository interface {
	Paginate(page int, size int, filter MovieFilter) (Page[Movie], error)
	FindByID(id uuid.UUID) (*Movie, error)
}
