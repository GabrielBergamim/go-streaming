package persistence

import (
	"github.com/example/go-streaming/api/domain/movie"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormMovieRepository struct {
	DB *gorm.DB
}

func NewGormMovieRepository(db *gorm.DB) *GormMovieRepository {
	return &GormMovieRepository{DB: db}
}

func (r *GormMovieRepository) Paginate(page int, size int,
	filter movie.MovieFilter) (movie.Page[movie.Movie], error) {
	var movies []movie.Movie
	scope := r.paginateScope(page, size)

	query := r.DB.Model(&movie.Movie{})
	if filter.Name != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Name+"%")
	}

	total := int64(0)
	if err := query.Count(&total).Error; err != nil {
		return movie.Page[movie.Movie]{}, err
	}

	if err := query.Scopes(scope).Find(&movies).Error; err != nil {
		return movie.Page[movie.Movie]{}, err
	}

	return movie.Page[movie.Movie]{
		Content: movies,
		TotalItems:  total,
		TotalPages: (total + int64(size) - 1) / int64(size),
		IsLast: (page+1)*size >= int(total),
	}, nil
}

func (r *GormMovieRepository) paginateScope(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

func (r *GormMovieRepository) FindByID(id uuid.UUID) (*movie.Movie, error) {
	var m movie.Movie
	if err := r.DB.First(&m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
