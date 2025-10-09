package movie

import (
	"fmt"
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

type MovieFilter struct {
	Name string `json:"name"`
}

type MovieDTO struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Genre       *string   `json:"genre,omitempty"`
	Director    *string   `json:"director,omitempty"`
	PosterURL   *string   `json:"poster_url,omitempty"`
	Rating      *float64  `json:"rating,omitempty"`
	Duration    *string   `json:"duration,omitempty"`
	ReleaseDate *string   `json:"release_date,omitempty"`
}

func ToMovieDTOs(movies []Movie) []MovieDTO {
	dtos := make([]MovieDTO, len(movies))

	for i, m := range movies {
		dtos[i] = ToMovieDTO(m)
	}

	return dtos
}

func ToMovieDTO (m Movie) MovieDTO {
	var durationStr *string
	if m.Duration != nil && m.Duration.Valid {
		d := IntervalToHHMMSS(*m.Duration)
		durationStr = &d
	}

	var releaseDateStr *string
	if m.ReleaseDate != nil {
		d := m.ReleaseDate.Format("2006-01-02")
		releaseDateStr = &d
	}

	return MovieDTO{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Genre:       m.Genre,
		Director:    m.Director,
		PosterURL:   m.PosterURL,
		Rating:      m.Rating,
		Duration:    durationStr,
		ReleaseDate: releaseDateStr,
	}
}

func IntervalToHHMMSS(iv pgtype.Interval) string {
	const (
		usPerSecond = int64(1_000_000)
		secPerMin   = int64(60)
		secPerHour  = int64(3600)
		secPerDay   = int64(86400)
	)

	totalSec := iv.Microseconds/usPerSecond + int64(iv.Days)*secPerDay
	neg := totalSec < 0
	if neg {
		totalSec = -totalSec
	}

	h := totalSec / secPerHour
	m := (totalSec % secPerHour) / secPerMin
	s := totalSec % secPerMin

	if neg {
		return fmt.Sprintf("-%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
