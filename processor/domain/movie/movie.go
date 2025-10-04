package movie

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Movie struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;column:id"`
	Title       string          `gorm:"type:varchar(255);not null;column:title"`
	Path        string          `gorm:"type:varchar(255);not null;column:path"`
	Description *string         `gorm:"type:text;column:description"`
	Genre       *string         `gorm:"type:varchar(255);column:genre"`
	Director    *string         `gorm:"type:varchar(255);column:director"`
	PosterURL   *string         `gorm:"type:text;column:poster_url"`
	Rating      *float64        `gorm:"type:double precision;column:rating"`
	Duration    *pgtype.Interval `gorm:"type:interval;column:duration"`
	ReleaseDate *time.Time       `gorm:"type:date;column:release_date"`
}

func (Movie) TableName() string {
	return "movies"
}
