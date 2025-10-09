package router

import (
	"os"

	"github.com/example/go-streaming/api/infra/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Router struct {
	app              *fiber.App
	moviesController *controller.MoviesController
}

func NewRouter(app *fiber.App, moviesController *controller.MoviesController) *Router {
	return &Router{
		app:              app,
		moviesController: moviesController,
	}
}

func (r *Router) SetUp() {
	r.app.Use(cors.New())
	r.staticFiles()
	r.moviesRoutes()
}

func (r *Router) moviesRoutes() {
	group := r.app.Group("/api/movies")
	group.Get("/", r.moviesController.GetMovies)
	group.Get("/:id", r.moviesController.GetMovie)
}

func (r *Router) staticFiles() {
	fiberStatic := fiber.Static{Compress: false, ByteRange: true}
	r.app.Static("/api/movie", os.Getenv("PUBLIC_FILES"), fiberStatic)
}
