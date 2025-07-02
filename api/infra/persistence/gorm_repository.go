package persistence

import (
	"github.com/example/go-streaming/api/domain/video"
	"gorm.io/gorm"
)

type GormVideoRepository struct {
	DB *gorm.DB
}

func NewGormVideoRepository(db *gorm.DB) *GormVideoRepository {
	return &GormVideoRepository{DB: db}
}

func (r *GormVideoRepository) FindByName(name string) ([]video.Video, error) {
	var videos []video.Video
	query := r.DB.Model(&video.Video{})
	if name != "" {
		query = query.Where("file_name ILIKE ?", "%"+name+"%")
	}
	if err := query.Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}
