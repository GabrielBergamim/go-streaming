package repository

import "github.com/example/go-streaming/processor/domain/movie"

type MovieRepository interface {
	Save(movie *movie.Movie) error
}
