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

func (r *GormVideoRepository) Paginate(page int, size int,
	filter video.VideoFilter) (video.Page[video.Video], error) {
	var videos []video.Video
	scope := r.paginateScope(page, size)

	query := r.DB.Model(&video.Video{})
	if filter.Name != "" {
		query = query.Where("file_name ILIKE ?", "%"+filter.Name+"%")
	}

	total := int64(0)
	if err := query.Count(&total).Error; err != nil {
		return video.Page[video.Video]{}, err
	}

	if err := query.Scopes(scope).Find(&videos).Error; err != nil {
		return video.Page[video.Video]{}, err
	}

	return video.Page[video.Video]{
		Content: videos,
		TotalItems:  total,
		TotalPages: (total + int64(size) - 1) / int64(size),
		IsLast: (page+1)*size >= int(total),
	}, nil
}

func (r *GormVideoRepository) paginateScope(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
