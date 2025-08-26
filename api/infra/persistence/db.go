package persistence

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(dsn string) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io.Writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // log slow queries
			LogLevel:                  logger.Info,            // Info shows SQL
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
			// Set to true to HIDE values (only placeholders) – leave false to see bound params
			ParameterizedQueries: false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
