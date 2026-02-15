package handlers

import (
	"cv_api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoHandlers struct {
	service *services.VideoService
}

func NewVideoHandlers(service *services.VideoService) *VideoHandlers {
	return &VideoHandlers{service: service}
}

func (h *VideoHandlers) GetVideos(c *gin.Context) {
	videos, err := h.service.GetVideos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch videos: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": videos,
	})
}
