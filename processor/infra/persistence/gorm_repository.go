package persistence

import (
	"github.com/example/go-streaming/processor/domain/repository"
	"github.com/example/go-streaming/processor/domain/video"
	"gorm.io/gorm"
)

type GormVideoRepository struct {
	DB *gorm.DB
}

func NewGormVideoRepository(db *gorm.DB) *GormVideoRepository {
	return &GormVideoRepository{DB: db}
}

func (r *GormVideoRepository) Save(v *video.Video) error {
	return r.DB.Create(v).Error
}

var _ repository.VideoRepository = (*GormVideoRepository)(nil)
