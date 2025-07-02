package repository

import "github.com/example/go-streaming/processor/domain/video"

type VideoRepository interface {
	Save(video *video.Video) error
}
