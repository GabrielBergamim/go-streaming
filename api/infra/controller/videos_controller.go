package controller

import (
	"strconv"

	"github.com/example/go-streaming/api/domain/video"
	"github.com/gofiber/fiber/v2"
)

type VideosController struct {
    VideoRepository video.VideoRepository
}

func NewVideosController(videoRepository video.VideoRepository) *VideosController {
	return &VideosController{
		VideoRepository: videoRepository,
	}
}

func (vc *VideosController) GetVideos(c *fiber.Ctx) error {
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

	pageable, err := vc.VideoRepository.Paginate(page,
		size, video.VideoFilter{Name: name})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch videos",
		})
	}

	return c.Status(fiber.StatusOK).JSON(pageable)
}
