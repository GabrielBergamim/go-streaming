package router

import (
	"os"

	"github.com/example/go-streaming/api/infra/controller"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app              *fiber.App
	videosController *controller.VideosController
}

func NewRouter(app *fiber.App, videosController *controller.VideosController) *Router {
	return &Router{
		app:              app,
		videosController: videosController,
	}
}

func (r *Router) SetUp() {
	r.staticFiles()
	r.videosRoutes()
}

func (r *Router) videosRoutes() {
	group := r.app.Group("/api/videos")
	group.Get("/", r.videosController.GetVideos)
}

func (r *Router) staticFiles() {
	fiberStatic := fiber.Static{Compress: false, ByteRange: true}
	r.app.Static("/api/video", os.Getenv("PUBLIC_FILES"), fiberStatic)
}
