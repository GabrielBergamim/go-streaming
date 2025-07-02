package video

type Repository interface {
	FindByName(name string) ([]Video, error)
}
