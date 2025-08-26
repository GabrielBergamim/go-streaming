package server

import (
	"log/slog"
	"mime"

	"github.com/example/go-streaming/api/infra/controller"
	"github.com/example/go-streaming/api/infra/persistence"
	"github.com/example/go-streaming/api/infra/router"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	Port string
	app  *fiber.App
	db   *gorm.DB
}

func NewServer(port string, dsn string) *Server {
	return &Server{
		Port: ":" + port,
		app: fiber.New(),
		db:  persistence.NewDB(dsn),
	}
}

func (s *Server) Start() error {
	s.setRoutes()

	mime.AddExtensionType(".m4s", "video/iso.segment")
	mime.AddExtensionType(".m3u8", "application/x-mpegURL")

	if err := s.app.Listen(s.Port); err != nil {
		slog.Error("Failed to start server", "error", err)
		return err
	}

	return nil
}

func (s *Server) setRoutes() {
	videoRepository := persistence.NewGormVideoRepository(s.db)
	router := router.NewRouter(s.app, controller.NewVideosController(videoRepository))
	router.SetUp()
}
