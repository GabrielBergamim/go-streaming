package video

type VideoRepository interface {
	Paginate(page int, size int, filter VideoFilter) (Page[Video], error)
}
