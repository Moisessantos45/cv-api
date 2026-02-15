package handlers

import (
	"cv_api/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProjectHandlers struct {
	service *services.ProjectService
}

func NewProjectHandlers(service *services.ProjectService) *ProjectHandlers {
	return &ProjectHandlers{service: service}
}

func (h *ProjectHandlers) GetProjects(c *gin.Context) {
	projects, err := h.service.GetProjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch projects: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": projects,
	})
}

func (h *ProjectHandlers) GetProjectById(c *gin.Context) {
	id := c.Param("id")
	if id == "" && strings.TrimSpace(id) == "" && len(strings.TrimSpace(id)) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project ID is required",
		})
		return
	}

	project, err := h.service.GetProject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch project: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": project,
	})
}

func (h *ProjectHandlers) UpdateCounter(c *gin.Context) {
	id := c.Param("id")
	if id == "" && strings.TrimSpace(id) == "" && len(strings.TrimSpace(id)) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project ID is required",
		})
		return
	}

	project, err := h.service.GetProject(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch project: " + err.Error(),
		})
		return
	}

	newProject, err := h.service.UpdateCounter(c.Request.Context(), id, project.CounterLikes+1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update counter: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": newProject,
	})
}
