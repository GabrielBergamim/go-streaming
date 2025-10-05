package movie

type MovieRepository interface {
	Paginate(page int, size int, filter MovieFilter) (Page[Movie], error)
}
