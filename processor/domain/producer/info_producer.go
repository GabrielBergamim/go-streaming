package producer

import "github.com/example/go-streaming/processor/domain/movie"

type InfoProducer interface {
	Send(mv *movie.Movie) error
}
