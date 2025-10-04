package persistence

import (
	"github.com/example/go-streaming/processor/domain/movie"
	"github.com/example/go-streaming/processor/domain/repository"
	"gorm.io/gorm"
)

type GormVideoRepository struct {
	DB *gorm.DB
}

func NewGormVideoRepository(db *gorm.DB) *GormVideoRepository {
	return &GormVideoRepository{DB: db}
}

func (r *GormVideoRepository) Save(v *movie.Movie) error {
	return r.DB.Create(v).Error
}

var _ repository.MovieRepository = (*GormVideoRepository)(nil)
