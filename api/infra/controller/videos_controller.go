package controller

import (
	"strconv"

	"github.com/example/go-streaming/api/domain/movie"
	"github.com/gofiber/fiber/v2"
)

type MoviesController struct {
    MovieRepository movie.MovieRepository
}

func NewMoviesController(movieRepository movie.MovieRepository) *MoviesController {
	return &MoviesController{
		MovieRepository: movieRepository,
	}
}

func (vc *MoviesController) GetMovies(c *fiber.Ctx) error {
	name := c.Query("name")
	page, err := strconv.Atoi(c.Query("page", "1"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page parameter",
		})
	}

	size, err := strconv.Atoi(c.Query("size", "10"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid size parameter",
		})
	}

	pageable, err := vc.MovieRepository.Paginate(page,
		size, movie.MovieFilter{Name: name})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch movies",
		})
	}

	pageableResponse := movie.Page[movie.MovieDTO] {
		Content:    movie.ToMovieDTOs(pageable.Content),
		TotalItems: pageable.TotalItems,
		TotalPages: pageable.TotalPages,
		IsLast:     pageable.IsLast,
	}

	return c.Status(fiber.StatusOK).JSON(pageableResponse)
}
