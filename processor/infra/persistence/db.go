package persistence

import (
	"github.com/example/go-streaming/processor/domain/movie"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&movie.Movie{}); err != nil {
		return nil, err
	}
	return db, nil
}
